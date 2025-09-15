"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const emitter_1 = require("./emitter");
const eventConsumer = new emitter_1.RabbitEvents();
eventConsumer.on("notification", async (content) => {
    console.log("JS Consumer: Got notification from emiiter");
    // some api function call to send email
    // perform typecasting for verifying data
    const u = content;
    await sendEmail(u);
});
async function sendEmail(u) {
    console.log("send email verification");
}
eventConsumer.start("break_prod", "amqp://localhost");
