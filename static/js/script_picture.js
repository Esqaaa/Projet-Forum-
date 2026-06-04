document.addEventListener("DOMContentLoaded", function () {
        const input = document.getElementById('image');
        const preview = document.getElementById('preview');
        const removeBtn = document.getElementById('removeImage');

        if (!input || !preview || !removeBtn) return;

        input.addEventListener('change', function () {
            const file = this.files[0];

            if (!file) {
                preview.style.display = 'none';
                removeBtn.style.display = 'none';
                preview.src = '';
                return;
            }

            if (!file.type.startsWith('image/')) {
                alert("Veuillez sélectionner une image valide.");
                this.value = "";
                return;
            }

            preview.src = URL.createObjectURL(file);
            preview.style.display = 'block';
            removeBtn.style.display = 'inline-flex'; 
        });

        removeBtn.addEventListener('click', function (e) {
            e.preventDefault();
            input.value = "";
            preview.src = "";
            preview.style.display = 'none';
            removeBtn.style.display = 'none';
        });
    });