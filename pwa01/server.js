import express from 'express'
import morgan from 'morgan';
import bcrypt from 'bcrypt';
import multer from 'multer';
import path from 'path';
import { fileURLToPath } from 'url';
import { checker, authentication }  from './database.js';


const app = express();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));
app.use(express.static(path.join(__dirname, 'public')));


const port = 3000;
const upload = multer()

const urlHome = 'http://localhost:3000/'

app.use(express.urlencoded({ extended: true}));
app.use(morgan('dev'));


//ROUTING
app.get('/', (req, res) => {
    res.status(200).render('index')
});

app.get('/sign-up', (req, res) => {
    res.status(200).render('sign-up');
})

app.get('/login', (req, res) => {
    res.status(200).render('login')
})


// sign in
app.post('/sign-up', upload.none(),    async (req, res) => {
    // destructure request
    const {email, password} = req.body;
    try {
        // generate salt-hashed password with bcrypt
        const salt = await bcrypt.genSalt();
        const hashedPassword = await bcrypt.hash(password, salt);

        // check database for existing user
        const result = await checker(email, hashedPassword);

        if (!result.success){
            res.status(409).json({message: 'User already exists'});
            return;
        }
        // return status code
        return res.status(201).json({message: 'Saved Successsfully'})

    } catch (err) {
        console.error(`MyServer Error: ${err}`)
    }

})

// handle login function 
app.post('/login', upload.none(), async (req, res) =>{
    const { email, password } = req.body;

    try {
        // use authentication to validate aganinst user credentials
        const result = await authentication(email, password);
        if (!result.success){
            return res.status(401).send('Incorrect Password/Email');
        } 

        return res.status(200).send('Logged In');


    } catch (err){
        console.log(`Error Within Server, ${err}`);
        res.status(500).send('Internal Server Error')
    }


})


// 404 Handling
app.use((req, res) => {
    res.status(404).render('404')
})

app.listen(port, ()=>{
    console.log(`Port ${port} active`)
    console.log(urlHome) 
})