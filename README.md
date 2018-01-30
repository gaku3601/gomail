# USAGE
## start

    docker build -t gomail .
    docker run -p 8888:8888 -e GMAIL=takuto.premium@gmail.com GMAILPW=xxxxxxxxxxxxxxx gomail

## mail post

    curl  -H 'Content-Type:application/json' http://localhost:8888/ -d '{"to": ["pro.gaku@gmail.com","t_hidaka@gmail.com"], "subject":"これはテストメールです", "message": "messageです。\r\nmessa"}'
