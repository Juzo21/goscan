document.addEventListener("DOMContentLoaded", function() {
    const sections = document.querySelectorAll("section");
    sections.forEach(section => {
        section.style.opacity = 0;
        section.style.transform = "translateY(20px)";
    });

    function revealSections() {
        sections.forEach((section, index) => {
            setTimeout(() => {
                section.style.opacity = 1;
                section.style.transform = "translateY(0)";
                section.style.transition = "all 0.6s ease";
            }, index * 200);
        });
    }

    revealSections();
});
