// Inicializar Animações
AOS.init({ duration: 1000, once: true });

// Menu Mobile
const hamburger = document.getElementById('hamburger');
const navMenu = document.getElementById('navMenu');

hamburger.addEventListener('click', () => {
    navMenu.classList.toggle('active');
    hamburger.classList.toggle('open');
});

// Efeito de contagem nos números (Stats)
const stats = document.querySelectorAll('.stat-number');
stats.forEach(stat => {
    const target = +stat.getAttribute('data-target');
    const updateCount = () => {
        const count = +stat.innerText;
        const speed = target / 100;
        if (count < target) {
            stat.innerText = Math.ceil(count + speed);
            setTimeout(updateCount, 20);
        } else {
            stat.innerText = target;
        }
    };
    updateCount();
});