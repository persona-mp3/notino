import { RabbitEvents } from "./emitter";

const eventConsumer = new RabbitEvents()
type User = {
  userName: string,
  email: string,
  firstName: string, 
  lastName: string,
}

eventConsumer.on("notification", async (content) => {
  console.log("JS Consumer: Got notification from emiiter")
  // some api function call to send email
  // perform typecasting for verifying data
  const u: User = content as User
  await sendEmail(u)
})

async function sendEmail(u: User){
  console.log("send email verification")
}
eventConsumer.start("break_prod",  "amqp://localhost")