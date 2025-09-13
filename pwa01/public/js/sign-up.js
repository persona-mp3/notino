const hamMenu = document.querySelector('.hamburger-menu');  
const offScreenMenu = document.querySelector('.off-screen-menu');
const signupButton = document.getElementById('signupBtn');

const signUpForm = document.getElementById('signup-form');
// const message = function(x){return `Thank you ${x}`};

async function sendData(formData, e){

    try {
        const res = await fetch('/sign-up', {
            method: 'POST',
            body: formData,
        });

        // alert user already exists
        if (res.status === 409){
            e.preventDefault();
            alert('User already Exists');
            return;
        }

        // alert user success
        if(res.status === 201){
            alert(`SUCCESS`);
            window.location.href = '/login'
        }

    } catch(err) {
        console.error(`Something occured during fetch event: ${err}`);
        alert(`SOMETHING WENT WRONG IN FETCH EVENT, ${err}`)
    }
}

hamMenu.addEventListener('click', () => {
    hamMenu.classList.toggle('active');
    offScreenMenu.classList.toggle('active')
})


signUpForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const formData = new FormData(signUpForm);
    sendData(formData, e);
    
})

