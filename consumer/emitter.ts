import amqp from "amqplib"
import { EventEmitter } from "events"

export class RabbitEvents extends EventEmitter {
  async start(queue: string, url: string) {
    try {
      const connection = await amqp.connect(url)
      const channel = await connection.createChannel()

      await channel.assertQueue(queue, { durable: false })
      console.log("-- Connection initialized")
      console.log("-- Channel created")
      console.log("-- Queue created")
      console.log("[o] Waiting for new messages\n")

      channel.consume(queue, (msg) => {
        if (msg) {
          console.log("[*] New Notification from Register\n")
          const content = msg.content.toString()
          console.log("\t ", content)
          channel.ack(msg)

          this.emit("notification", content)
        }
      })
    } catch (err) {
      console.error("kaboom: error ocuured")
      throw err
    }
  }
}
