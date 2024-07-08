document.getElementById('paymentForm').addEventListener('submit', async function (event) {
    event.preventDefault();

    const formData = {
        userId: parseInt(document.getElementById('userId').value),
        amount: parseFloat(document.getElementById('amount').value),
        method: document.getElementById('method').value,
        transactionId: document.getElementById('transactionId').value
    };

    try {
        const response = await fetch('/payments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            document.getElementById('message').textContent = 'Payment successful!';
            document.getElementById('message').classList.add('alert', 'alert-success');
        } else {
            const errorText = await response.text();
            document.getElementById('message').textContent = 'Payment failed: ' + errorText;
            document.getElementById('message').classList.add('alert', 'alert-danger');
        }
    } catch (error) {
        document.getElementById('message').textContent = 'Payment failed: ' + error.message;
        document.getElementById('message').classList.add('alert', 'alert-danger');
    }
});
document.getElementById('paymentForm').addEventListener('submit', async function (event) {
    event.preventDefault();

    const formData = {
        userId: parseInt(document.getElementById('userId').value),
        amount: parseFloat(document.getElementById('amount').value),
        method: document.getElementById('method').value,
        transactionId: document.getElementById('transactionId').value
    };

    try {
        const response = await fetch('/payments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            document.getElementById('message').textContent = 'Payment successful!';
            document.getElementById('message').classList.add('alert', 'alert-success');
        } else {
            const errorText = await response.text();
            document.getElementById('message').textContent = 'Payment failed: ' + errorText;
            document.getElementById('message').classList.add('alert', 'alert-danger');
        }
    } catch (error) {
        document.getElementById('message').textContent = 'Payment failed: ' + error.message;
        document.getElementById('message').classList.add('alert', 'alert-danger');
    }
});
