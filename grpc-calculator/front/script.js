function submitExpression() {
    const expressionInput = document.getElementById("expression");
    const expression = expressionInput.value;

    const data = {
        expression: expression,
    };

    fetch('http://localhost:8081/culc', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getTokenFromCookie()}`
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            expressionInput.value = "";
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        // After successful request execution, display the result on the page
        const resultElement = document.createElement("div");
        resultElement.textContent = `${expression} = ${data.result}`;

        const calcForm = document.getElementById("calc-form");
        calcForm.appendChild(resultElement);

        expressionInput.value = "";
    })
    .catch(error => {
        console.error('Calculation error:', error);
        const resultElement = document.createElement("div");
        const calcForm = document.getElementById("calc-form");
        resultElement.textContent = `${expression} = ${'calculation error'}`;
        calcForm.appendChild(resultElement);

        expressionInput.value = "";
    });
}

function getTokenFromCookie() {
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        const [name, value] = cookie.trim().split('=');
        if (name === 'token') {
            return value;
        }
    }
    return null;
}

function GetStatus() {
    fetch('http://localhost:8081/allclient', {
        method: 'GET',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        displayClientStatus(data); // Display client information
    })
    .catch(error => {
        console.error('Login error:', error);
        alert('Login error');
    });
}

function displayClientStatus(data) {
    // Clear the container before updating
    const statusContainer = document.getElementById('status-container');
    statusContainer.innerHTML = '';

    // Check for the presence of the "statuses" key in the data object
    if (data.hasOwnProperty('statuses')) {
        const statuses = data.statuses;

        // Loop through the statuses array
        statuses.forEach(client => {
            const clientInfo = document.createElement('p');

            // Display client status information
            clientInfo.textContent = `Server ID: ${client.serverId}, Active: ${client.active}, Total: ${client.all}`;

            statusContainer.appendChild(clientInfo);
        });
    } else {
        // If the "statuses" key is absent, display the corresponding message
        const clientInfo = document.createElement('p');
        clientInfo.textContent = `No client data`;
        statusContainer.appendChild(clientInfo);
    }
}

function GetHistory(){
    fetch('http://localhost:8081/history', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${getTokenFromCookie()}`
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        displayExpressionHistory(data);
    })
    .catch(error => {
        console.error('Login error:', error);
        alert('Login error');
    });

}

function displayExpressionHistory(data) {
    // Clear the container before updating
    const historyContainer = document.getElementById('history-container');
    historyContainer.innerHTML = '';

    // Check for the presence of the "expressions" key in the data object
    if (data.hasOwnProperty('expressions')) {
        const expressions = data.expressions;

        // Loop through the expressions array
        expressions.forEach(expressionObj => {
            const expressionInfo = document.createElement('p');

            // Display expression information
            expressionInfo.textContent = `ID: ${expressionObj.id}, Expression: ${expressionObj.expres}`;

            // Add expression information to the container
            historyContainer.appendChild(expressionInfo);

            // Add separator
            const separator = document.createElement('hr');
            historyContainer.appendChild(separator);
        });
    } else {
        // If the "expressions" key is absent, display the corresponding message
        const noHistoryMessage = document.createElement('p');
        noHistoryMessage.textContent = `No expression history`;
        historyContainer.appendChild(noHistoryMessage);
    }
}

function SetConfig(){
    const Sum = document.getElementById('sum').value;
    const div = document.getElementById('div').value;
    const Ex = document.getElementById('exp').value;
    const Min = document.getElementById('min').value;
    const Multp = document.getElementById('mult').value;
  
    const data = {
        Div: div,
        Exponent: Ex,
        Minus: Min,
        MultP: Multp,
        Plus: Sum
    };
   
    fetch('http://localhost:8081/config', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        Sum.value="";
        div.value="";
        Ex.value ="";
        Min.value="";
        Multp.value="";
        return response.json();
    })
    .then(data => {
        console.log('Successful login:', data);
        alert('Successful login');
    })
    .catch(error => {
        console.error('Login error:', error);
        alert('Login error');
    });
}

function openModal() {
    document.getElementById('modal').style.display = 'block';
    // Add click event handler to the entire document
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
}

function addClient() {
    const countworkerInput = document.getElementById('numberOfGoroutines').value;
    const data = {
        countworker: +countworkerInput
    };
    console.log(data)
    fetch('http://localhost:8081/newclient', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log('Successful adding client:', data);
        alert('Successful adding client');
    })
    .catch(error => {
        console.error('Error adding client:', error);
        alert('Error adding client');
    });
}

