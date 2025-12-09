function showErrorToast(text) {
    Toastify({
        text: text,
        duration: 3000,
        gravity: "bottom",
        position: "center",
        style: {
            background: "linear-gradient(to right, #ad2d0dff, #901b09ff)",
            fontFamily: 'Arial, sans-serif',
        },
    }).showToast();
}

function showSuccessToast(text) {
    Toastify({
        text: text,
        duration: 2000,
        gravity: "bottom",
        position: "center",
        style: {
            background: "linear-gradient(to right, #2dad0dff, #05630bff)",
            fontFamily: 'Arial, sans-serif',
        },
    }).showToast();
}