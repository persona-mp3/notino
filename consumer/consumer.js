"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const emitter_1 = require("./emitter");
const mail_1 = __importDefault(require("@sendgrid/mail"));
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
const API_KEY = process.env.SG_KEY;
if (!API_KEY) {
    throw ("Could not read api key");
}
mail_1.default.setApiKey(API_KEY);
const eventConsumer = new emitter_1.RabbitEvents();
const msgFmt = (name) => `
<h2> Hey there ${name} </h2>

<p>You have just signed up to create an account with ImportantBusiness.</p>
<p>
  For your first task, please visit <a>https://github.com/persona-mp3/proto.git</a>.
  When you've been able to complete the task, don't forget to send the your answers to 
  us through the website.
</p> 

  <p>persona - 20 Verdant Street, The Bikini Bottom</p>
`;
eventConsumer.on("notification", async (content) => {
    console.log("JS Consumer: Got notification from emiiter");
    const parsedContent = JSON.parse(content);
    console.log(parsedContent);
    const u = parsedContent;
    try {
        await sendEmail(u);
    }
    catch (err) {
        console.log(err);
    }
});
async function sendEmail(u) {
    console.log("sending email verification to -> ", u.email);
    const msg = {
        to: u.email,
        from: "personacodes@gmail.com",
        subject: "Account Confirmation",
        html: msgFmt(`${u.firstName} ${u.lastName}`)
    };
    try {
        const res = await mail_1.default.send(msg);
        console.log("sent email, reading response headers");
        const statusCode = res[0].statusCode;
        const headers = res[0].headers;
        console.log("status code returned -> %s", statusCode);
    }
    catch (err) {
        throw err;
    }
}
eventConsumer.start("break_prod", "amqp://localhost");
