function setUsernameOnFooter() {
    const username = JSON.parse(document.getElementById("profile-selected-username").textContent);
    document.getElementById('profile-selected').innerText = ": " + username;
}

let previousSelectedItem = null;

function filterItems() {
    const input = document.getElementById('filter-items-input');
    const filter = input.value.toUpperCase().trim();
    const itemList = document.getElementById("item-list");
    const li = itemList.getElementsByTagName('li');

    // Loop through all list items, and hide those who don't match the search query
    for (let i = 0; i < li.length; i++) {
        const txtValue = (li[i].textContent || li[i].innerText).toUpperCase().trim();
        const itemId = li[i].getAttribute('data-item-id').toUpperCase().trim();
        const itemCategory = li[i].parentElement.getAttribute('data-category').toUpperCase().trim();
        if (txtValue.indexOf(filter) > -1 || itemCategory.indexOf(filter) > -1 || itemId.indexOf(filter) > -1) {
            li[i].style.display = "";
        } else {
            li[i].style.display = "none";
        }
    }

    // hide empty categories
    const section = itemList.getElementsByTagName('details');
    for (let i = 0; i < section.length; i++) {
        const innerLis = section[i].getElementsByTagName('li');
        let anyLiVisible = false;
        for (let j = 0; j < innerLis.length; j++) {
            if (innerLis[j].style.display === "") {
                anyLiVisible = true;
                break;
            }
        }
        if (anyLiVisible) {
            section[i].style.display = "";
        } else {
            section[i].style.display = "none";
        }
    }
}

function selectItem(element) {
    const classToToggle = 'text-primary';
    if (previousSelectedItem) {
        previousSelectedItem.classList.remove(classToToggle);
    }
    element.classList.add(classToToggle);
    previousSelectedItem = element;
}

function selectItemFromKeyboard(event, element) {
    if (event.key === 'Enter') {
        selectItem(element);
    }
}

function filterUserWeapons() {
    const input = document.getElementById('filter-user-weapons-input');
    const filter = input.value.toUpperCase().trim();
    const itemList = document.getElementById("weapons-list");
    const cards = itemList.getElementsByClassName('card-to-filter');

    // Loop through all list items, and hide those who don't match the search query
    for (let i = 0; i < cards.length; i++) {
        const title = cards[i].getElementsByTagName('h2')[0];
        const txtValue = (title.textContent || title.innerText).toUpperCase().trim();
        if (txtValue.indexOf(filter) > -1) {
            cards[i].style.display = "";
        } else {
            cards[i].style.display = "none";
        }
    }
}

function showModal(event) {
    const dialog = document.getElementById(event.getAttribute('data-dialog-target'));
    dialog.showModal();
}

function filterMagazineLoadout() {
    const input = document.getElementById('filter-magazine-loadout-input');
    const filter = input.value.toUpperCase().trim();
    const itemList = document.getElementById("magazine-loadout-list");
    const cards = itemList.getElementsByClassName('card-to-filter');

    // Loop through all list items, and hide those who don't match the search query
    for (let i = 0; i < cards.length; i++) {
        const title = cards[i].getElementsByTagName('h2')[0];
        const txtValue = (title.textContent || title.innerText).toUpperCase().trim();
        if (txtValue.indexOf(filter) > -1) {
            cards[i].style.display = "";
        } else {
            cards[i].style.display = "none";
        }
    }
}

window.setUsernameOnFooter = setUsernameOnFooter;
window.filterItems = filterItems;
window.selectItem = selectItem;
window.selectItemFromKeyboard = selectItemFromKeyboard;
window.filterUserWeapons = filterUserWeapons;
window.showModal = showModal;
window.filterMagazineLoadout = filterMagazineLoadout;

let previousToast = null;

htmx.on("showAddItemMessage", (e) => {
    if (previousToast) {
        clearTimeout(previousToast)
    }
    const toastElement = document.getElementById("success-toast")
    const toastBody = document.getElementById("success-toast-message")
    toastBody.innerText = e.detail.value;
    toastElement.classList.remove("hidden")
    previousToast = setTimeout(() => {
        toastElement.classList.add("hidden")
    }, 2000)
});

window.winterEvent = function () {
    const elem = document.getElementById("winter-event");
    const today = new Date();
    if (elem && today.getMonth() === 11 && today.getDate() >= 24) {
        elem.innerHTML = `
       <div class="snowflakes" aria-hidden="true">
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        <div class="snowflake">
            <div class="inner">❅</div>
        </div>
        </div>
       `;
    }

}

