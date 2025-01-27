let currentSlide = 0;

function moveCarousel(direction) {
    const container = document.querySelector('.carousel-container');
    const items = document.querySelectorAll('.accordion-item');
    const totalSlides = items.length;
    
    // Update current slide
    currentSlide = Math.max(0, Math.min(currentSlide + direction, totalSlides - 1));
    
    // Move the carousel
    container.style.transform = `translateX(-${currentSlide * 100}%)`;
    
    // Update button states
    const prevButton = document.querySelector('.carousel-nav.prev');
    const nextButton = document.querySelector('.carousel-nav.next');
    
    prevButton.disabled = currentSlide === 0;
    nextButton.disabled = currentSlide === totalSlides - 1;
    
    // Trigger htmx to load charts if needed
    const currentItem = items[currentSlide];
    const charts = currentItem.querySelectorAll('[hx-trigger="load"]');
    charts.forEach(chart => {
        if (!chart.innerHTML.includes('svg')) {
            htmx.trigger(chart, 'load');
        }
    });
}