from operator import truediv
from flask import Flask, request, jsonify
from twilio.rest import Client
import os
import json

app = Flask(__name__)
client = Client(os.getenv("AccountSid"), os.getenv("AuthToken"))

#def create_app():
#   return app

@app.route("/",methods = ['POST', 'GET'])
def ping():
    return '<h1>Hello, World!</h1>'

@app.route("/sms",methods = ['POST'])
def sms_request():
    if request.is_json:
        data = request.get_json()
        message = client.messages.create(
            # TODO: Better content in SMS. 
                            body=json.dumps(data),
                            # TODO: Alphanumeric SenderID to display system name.
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
        data = request.get_json()
        call = client.calls.create(
            # TODO: Better content in call. 
                        twiml='<Response><Say>' + json.dumps(data) + '</Say></Response>',
                        to=request.args.get('receiver'),
                        from_=os.getenv("Sender")
                    )
        print(call.sid)
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

if __name__=="__main__":
    app.run