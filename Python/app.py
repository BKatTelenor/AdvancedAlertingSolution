from flask import Flask, request, jsonify
from twilio.rest import Client
import os

app = Flask(__name__)
client = Client(os.getenv("AccountSid"), os.getenv("AuthToken"))

@app.route("/",methods = ['POST', 'GET'])
def ping():
    return '<h1>Hello, World!</h1>'

@app.route("/sms",methods = ['POST'])
def sms_request():
    if request.is_json:
        content = "Hei Bjørn Kristian! Din Server er nede! du Suger!"
        message = client.messages.create(
                            body=content,
                            from_=os.getenv("Sender"),
                            to=request.args.get('receiver')
                        )
        print(message.sid)
        return '<h1>SMS</h1>'
    else:
        return

@app.route("/call",methods = ['POST'])
def call_request():
    if request.is_json:
        content = "Hei Bjørn Kristian! Din Server er nede! du Suger!Hei Bjørn Kristian! Din Server er nede! du Suger!"
        call = client.calls.create(
                        twiml='<Response><Say>' + content + '</Say></Response>',
                        to=request.args.get('receiver'),
                        from_=os.getenv("Sender")
                    )
        print(call.sid, "\n")
        return '<h1>Call</h1>'
    else:
        return

@app.route("/callandsms",methods = ['POST'])
def call_and_sms():
    if request.is_json:
        sms_request()
        call_request()
        return '<h1>Call and SMS</h1>'
    else:
        return