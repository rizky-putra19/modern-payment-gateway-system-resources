package email

const EmailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Hypay</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; color: #333; }
        .email-container { width: 100%; max-width: 600px; background-color: #fff; margin: 20px auto; padding: 20px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
        .header { background-color: #8A2BE2; color: #ffffff; padding: 20px; text-align: center; border-radius: 10px 10px; }
        .content { padding: 20px; text-align: left; line-height: 1.5; }
        .footer { text-align: center; padding: 10px 20px; font-size: 12px; color: #999; }
        .info { background-color: #eee; padding: 10px; margin: 10px 0; border-radius: 10px 10px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <img src="https://res.cloudinary.com/ddtewkcqc/image/upload/v1723206638/lwi2ojljn7ozab8l0siq.png" alt="Hypay Logo" style="max-width: 100px;">
            <h1>Welcome to Hypay!</h1>
        </div>
        <div class="content">
            <p>Dear Merchant,</p>
            <p>You are now registered with Hypay. Here are your credentials to access our dashboard:</p>
            <div class="info">
                <p><strong>Username:</strong> {{.Username}}</p>
                <p><strong>Password:</strong> {{.Password}}</p>
                <p><strong>Pin:</strong> {{.Pin}}</p>
            </div>
            <p>Please ensure that you keep this information secure at all times.</p>
            <p>Please change all credentials like password and pin on your manage profile menu on our merchant dashboard.</p>
        </div>
        <div class="footer">
            © Hypay Indonesia. All rights reserved.
        </div>
    </div>
</body>
</html>
`
