package authgRPC

import (
	"context"
	"culc/iternal/config"
	"culc/iternal/lib/other"
	"culc/iternal/lib/postfix"
	"culc/iternal/model"
	"culc/iternal/services/worker"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	culcv1 "github.com/ragnack97/protoculc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	Register(ctx context.Context, email string, password string) (uint, error)
	SaveExpression(ctx context.Context, exp string, uid float64) (int64, []string, error)
	UpdateSubEx(ctx context.Context, idex int64, userid float64, subEx []string) error
	GetExHistory(ctx context.Context, userid float64) ([]model.Expression, error)
	UpdateExTimeReady(ctx context.Context, idex int64, userid float64) error
	GetUnREadyEx() ([]model.Expression, error)
	ErrorinEx(float64, int64) error
}
type ServerAPI struct {
	culcv1.UnimplementedAuthServer
	auth         Auth
	CalcServ     *worker.CalculatorServer
	shutdownOnce sync.Once
}

func Register(gRPC *grpc.Server, auth Auth, cserv *worker.CalculatorServer) {
	culcv1.RegisterAuthServer(gRPC, &ServerAPI{auth: auth, CalcServ: cserv})

}
func (s *ServerAPI) Login(ctx context.Context, req *culcv1.LoginReq) (*culcv1.LoginRes, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty email or password")
	}
	s.shutdownOnce.Do(func() {
		go s.ChekShutDownEx()
	})
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}
	return &culcv1.LoginRes{
		Token: token,
	}, nil
}
func (s *ServerAPI) Register(ctx context.Context, req *culcv1.RegisterReq) (*culcv1.RegisterRes, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty email or password")
	}

	UserID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("error in register")
	}
	return &culcv1.RegisterRes{
		UserId: int64(UserID),
	}, nil
}
func (s *ServerAPI) Calculate(ctx context.Context, req *culcv1.CalculateReq) (*culcv1.CalculateRes, error) {
	if req.GetExpression() == "" {
		return nil, fmt.Errorf("empty expression")

	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, fmt.Errorf("missing token")
	}
	token := strings.TrimPrefix(authHeaders[0], "Bearer ")

	uid, timeexp, err := other.ParseToken(token)
	if err != nil {
		return nil, err
	}

	ok = time.Now().After(timeexp)
	if ok {
		return nil, fmt.Errorf("token просрочен")
	}

	idex, subEx, err := s.auth.SaveExpression(ctx, req.GetExpression(), uid)
	if err != nil {
		return nil, err
	}
	res, err := s.GetRes(ctx, uid, idex, subEx)
	if err != nil {
		return nil, fmt.Errorf("err: %e", err)
	}

	return &culcv1.CalculateRes{Result: res}, nil
}
func (s *ServerAPI) GetRes(ctx context.Context, uid float64, idex int64, subEx []string) (string, error) {
	wg := sync.WaitGroup{}

	for len(subEx) != 1 {
		errCh := make(chan error, len(subEx))
		results := make(map[int]string)

		for i := 0; i < len(subEx)-2; i++ {
			if isNumeric(subEx[i]) && isNumeric(subEx[i+1]) && strings.ContainsAny(subEx[i+2], "+-/*^") {

				clientId := s.CalcServ.GetClientWithMaxFreeWorkers()
				wg.Add(1)
				localI := i
				go func(clientid string, index int, results map[int]string) {
					defer wg.Done()

					client := s.CalcServ.GetClientByID(clientid)
					res, err := client.CountEx(subEx[i], subEx[i+1], subEx[i+2])

					if err != nil {

						errCh <- fmt.Errorf(" %v", err)

						return
					}

					results[index] = fmt.Sprint(res)

				}(clientId, localI, results)
			}

		}

		wg.Wait()
		close(errCh)

		var errs []error
		for err := range errCh {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			err := s.auth.ErrorinEx(uid, idex)
			if err != nil {

			}
			return "", fmt.Errorf("calculation error: %v", errs)
		}
		for idx, val := range results {

			subEx = append(append(subEx[:idx], val), subEx[idx+1:]...)
			subEx = append(append(subEx[:idx+1], " "), subEx[idx+2:]...)
			subEx = append(append(subEx[:idx+2], " "), subEx[idx+3:]...)

		}
		subEx = other.RemoveEmptyStrings(subEx)

		err := s.auth.UpdateSubEx(ctx, idex, uid, subEx)
		if err != nil {
			fmt.Println("error in update")
		}
	}
	s.auth.UpdateExTimeReady(ctx, idex, uid)
	return subEx[0], nil
}
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
func (s *ServerAPI) NewClient(ctx context.Context, req *culcv1.ClientReq) (*culcv1.ClientRes, error) {
	rand.Seed(time.Now().UnixNano())

	if req.GetCountworker() == 0 {
		return nil, fmt.Errorf("количество горутин==0")
	}
	randomNumber := rand.Intn(100)
	s.CalcServ.AddClient(fmt.Sprintf("%d", randomNumber), int(req.GetCountworker()))

	return &culcv1.ClientRes{Res: fmt.Sprint(randomNumber)}, nil

}
func (s *ServerAPI) StreamServerStatuses(ctx context.Context, req *culcv1.StreamServerStatusesRequest) (*culcv1.StreamServerStatusesResponse, error) {

	statuses := make([]*culcv1.ServerStatus, 0)
	for clientId, client := range s.CalcServ.GetFreeWorkerCounts() {
		statuses = append(statuses, &culcv1.ServerStatus{
			ServerId: clientId,
			Active:   int32(client.ActiveWorkers),
			All:      int32(client.AllWorkers),
		})
	}

	// Отправляем обновленные статусы клиенту
	return &culcv1.StreamServerStatusesResponse{
		Statuses: statuses,
	}, nil

}
func (s *ServerAPI) GetHistoryEx(ctx context.Context, req *culcv1.HistoryReq) (*culcv1.HistoryRes, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, fmt.Errorf("missing token")
	}
	token := strings.TrimPrefix(authHeaders[0], "Bearer ")

	uid, timeexp, err := other.ParseToken(token)
	if err != nil {
		return nil, err
	}

	ok = time.Now().After(timeexp)
	if ok {
		return nil, fmt.Errorf("token просрочен")
	}
	Ex, err := s.auth.GetExHistory(ctx, uid)
	if err != nil {
		return nil, err
	}
	var protoExpressions []*culcv1.Expression
	for _, exp := range Ex {
		m := ""
		if len(exp.SubEx) > 0 {
			m = exp.SubEx[:len(exp.SubEx)-1]
		}
		protoExpressions = append(protoExpressions, &culcv1.Expression{
			Id:     int32(exp.ID),
			Expres: exp.EX + "=" + m,
		})
	}

	return &culcv1.HistoryRes{Expressions: protoExpressions}, nil
}
func (s *ServerAPI) UpdateConfig(ctx context.Context, in *culcv1.ConfigReq) (*culcv1.ConfigRes, error) {
	var NewCFG config.ConfigTime

	NewCFG.Division = time.Duration(in.GetDiv())
	NewCFG.Exponent = time.Duration(in.GetExponent())
	NewCFG.Minus = time.Duration(in.Minus)
	NewCFG.MultiP = time.Duration(in.GetMultP())
	NewCFG.Plus = time.Duration(in.GetPlus())
	s.CalcServ.UpdateTime(NewCFG)
	fmt.Println(NewCFG.Division)
	return &culcv1.ConfigRes{}, nil
}
func (s *ServerAPI) ChekShutDownEx() {
	for {

		cfg := s.CalcServ.GetClientByID(s.CalcServ.GetClientWithMaxFreeWorkers()).GetConfig()
		maxDuration := cfg.Plus

		if cfg.Minus > maxDuration {
			maxDuration = cfg.Minus

		}
		if cfg.Division > maxDuration {
			maxDuration = cfg.Division

		}
		if cfg.MultiP > maxDuration {
			maxDuration = cfg.MultiP

		}
		if cfg.Exponent > maxDuration {
			maxDuration = cfg.Exponent

		}

		ex, err := s.auth.GetUnREadyEx()
		if err != nil {
			log.Print("ошибка в получении упавших выраженний")
		}
		ttl := time.Now().Add(maxDuration * 5)
		for _, exxpres := range ex {

			if ttl.After(exxpres.TimeUpdate) {
				go func() {

					sb := postfix.ParsSlice(exxpres.SubEx)

					s.GetRes(context.Background(), float64(exxpres.UserId), int64(exxpres.ID), sb)
				}()
			}

		}
		time.Sleep(time.Minute * 5)
	}
}
