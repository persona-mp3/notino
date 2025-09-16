"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.RabbitEvents = void 0;
const amqplib_1 = __importDefault(require("amqplib"));
const events_1 = require("events");
class RabbitEvents extends events_1.EventEmitter {
    async start(queue, url) {
        try {
            const connection = await amqplib_1.default.connect(url);
            const channel = await connection.createChannel();
            await channel.assertQueue(queue, { durable: false });
            console.log("-- Connection initialized");
            console.log("-- Channel created");
            console.log("-- Queue created");
            console.log("[o] Waiting for new messages\n");
            channel.consume(queue, (msg) => {
                if (msg) {
                    console.log("[*] New Notification from Register\n");
                    const content = msg.content.toString();
                    console.log("\t ", content);
                    channel.ack(msg);
                    this.emit("notification", content);
                }
            });
        }
        catch (err) {
            console.error("kaboom: error ocuured");
            throw err;
        }
    }
}
exports.RabbitEvents = RabbitEvents;
