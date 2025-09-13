if ('serviceWorker' in navigator){
    console.log('Service Worker Supported')

    window.addEventListener('load', () => {
        //access sw object
        navigator.serviceWorker
            .register('../sw.js')
            .then((result) => console.log('SW registered Successfully'))
            .catch((err) => console.log(`Error occured, ${err}`))
    });
} else{
    console.log('SW not supported')
}


const hamMenu = document.querySelector('.hamburger-menu');      
const offScreenMenu = document.querySelector('.off-screen-menu');
const signupButton = document.getElementById('signupBtn');

let viewPort = 800;

if (window.innerWidth > viewPort){
    gsap.from('.nav-container', {duration: 1.5, x:'-100%', ease: 'expo'});
    gsap.registerPlugin(ScrollTrigger);
    gsap.from('.library-section', 
        {
        scrollTrigger: {
            trigger: '.library-section',
            toggleActions: 'restart none none none',
        },

        duration: 2.5,
        x: '400%',
        ease: 'expo',
    });
}

hamMenu.addEventListener('click', () => {
    hamMenu.classList.toggle('active');
    offScreenMenu.classList.toggle('active')
})
