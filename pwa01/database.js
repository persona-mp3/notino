import mysql from 'mysql2';
import dotenv from 'dotenv';
import bcrypt from 'bcrypt'
dotenv.config()

// create pool connection
const pool = mysql.createPool( {
    host: process.env.MYSQL_HOST,
    user: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DATABASE,
}).promise();

// save user details
export async function saveDetails(email, password){
    const [user] = await pool.query(`
        INSERT INTO Users(email, password)
        VALUES (?, ?)
        `, [email, password])

    return user
};

// check for existing users
async function existingUser(email){
    const [checkUser]=  await pool.query(`
        SELECT email FROM Users
        WHERE email = (?)
        `, [email]);
    
    const existingUser = checkUser[0] 
    
    // user exists because database returns email
    if (existingUser !== undefined){        
        return true;
    }
    
    // user doesn't exist and returns undefined
    return false;

}

// call stack 
export async function checker(email, password){
    try{
        const userExists = await existingUser(email);

        // returns false if user exists
        if (userExists){
            return {success: false, message: 'USER EXISTS'}
        };

        // saves user to database if user doesn't exist
        await saveDetails(email, password);
        
        return {
            success: true, 
            message: 'USER CREATED'
        }
    } catch (err){
        console.error(`Duplicate Entry, User Already exists ${err}`);
        return false
    }
};

// get password based of email
async function checkPassword(email){
    const [isEmail] = await pool.query(`
        SELECT password FROM Users
        WHERE email = (?)
        `, [email])
    
    return isEmail[0].password
}

// call stack for login functiona;ity
export async function authentication(email, userPassword){

    try {
        // get password if user exists
        const userExists = await existingUser(email);
        if (!userExists){
            return {success: false, message: 'User does not exist'};
        }


        // compare password from user and stored password (hashed)
        const hashedPassword = await checkPassword(email);
        const match = await bcrypt.compare(userPassword, hashedPassword);

        // returns false if the passwords don't match
        if(!match){
            console.log('Incorrect Credentials. Try Again');
            return false;
        }
        // returns true if password match
        if(match){
            console.log('Logged IN');
            return {success: true};
    
        }
    } catch (err) {
        console.log(`Error occured during Authentication: ${err}`)
    }
}

