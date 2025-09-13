const loginForm = document.getElementById('login-form');
const hamMenu = document.querySelector('.hamburger-menu');  
const offScreenMenu = document.querySelector('.off-screen-menu');
const signupButton = document.getElementById('signupBtn');

// fetchAPI to send data
async function validateCredentials(formData, e){
    try{
        const res = await fetch('/login', 
            {
                method: 'POST', 
                body: formData
            })
        
        // if login credetials are invalid
        if (res.status === 401){
            e.preventDefault();
            alert('Invalid Creds');
            return;
        }

        if (res.status === 200){
            alert(`LOGGED IN`)
        }
        
    } catch (err){
        alert('Something went wrong');
        console.error(`My err: ${err}`)
    }
};

// trigger menu animation
hamMenu.addEventListener('click', () => {

    hamMenu.classList.toggle('active');
    offScreenMenu.classList.toggle('active')
})

loginForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const formData = new FormData(loginForm);
    validateCredentials(formData, e);
});

