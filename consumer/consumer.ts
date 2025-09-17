import { RabbitEvents } from "./emitter";
import sgMail from "@sendgrid/mail"
import dotenv from "dotenv"

dotenv.config()
const API_KEY = process.env.SG_KEY 
if (!API_KEY) {
  throw ("Could not read api key")
}
sgMail.setApiKey(API_KEY)

const eventConsumer = new RabbitEvents()
type User = {
  id: string,
  userName: string,
  email: string,
  firstName: string,
  lastName: string,
}

// At the moment, all emails sent usually go to a users spam folder because 
// of certain verfication issues, and since I don't have a valid domain yet 
const msgFmt = (name: any) => `
<h2> Hey there ${name} </h2>

<p>You have just signed up to create an account with ImportantBusiness.</p>
<p>
  For your first task, please visit <a>https://github.com/persona-mp3/proto.git</a>.
  When you've been able to complete the task, don't forget to send the your answers to 
  us through the website.
</p> 

  <p>persona - 20 Verdant Street, The Bikini Bottom</p>
`
eventConsumer.on("notification", async (content) => {
  console.log("JS Consumer: Got notification from emiiter")
  const parsedContent = JSON.parse(content)
  console.log(parsedContent)
  const u: User = parsedContent as User
  try {
    await sendEmail(u)
  } catch (err) {
    console.log(err)
  }
})



async function sendEmail(u: User) {
  console.log("sending email verification to -> ", u.email)
  const msg = {
    to: u.email,
    from: "personacodes@gmail.com",
    subject: "Account Confirmation",
    html: msgFmt(`${u.firstName} ${u.lastName}`)
  }

  try {
    const res = await sgMail.send(msg)
    console.log("sent email, reading response headers")
    const statusCode = res[0].statusCode

    console.log("status code returned -> %s", statusCode)
  } catch (err) {
    throw err
  }
}
eventConsumer.start("create_user", "amqp://localhost")
