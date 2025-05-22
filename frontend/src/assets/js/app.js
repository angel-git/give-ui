let previousSelectedItem = null;
let previousSelectedkit = null;
let previousToast = null;
let previousToastElement = null;

function setUsernameOnFooter() {
    const username = JSON.parse(document.getElementById("profile-selected-username").textContent);
    const el = document.getElementById('profile-selected')
    el.innerHTML = `<span>: ${username}</span><button class="btn btn-xs btn-link" hx-get="/reload-profiles" hx-target="#main" hx-swap="innerHTML" hx-disabled-elt="this">Switch profile</button>`;
    window.htmx.process(el);
}

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

function selectKit(element) {
    const classToToggle = 'text-primary';
    if (previousSelectedkit) {
        previousSelectedkit.classList.remove(classToToggle);
    }
    element.classList.add(classToToggle);
    previousSelectedkit = element;
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

function filterKits() {
    const input = document.getElementById('filter-kits-input');
    const filter = input.value.toUpperCase().trim();
    const itemList = document.getElementById("kits-list");
    const li = itemList.getElementsByTagName('li');


    // Loop through all list items, and hide those who don't match the search query
    for (let i = 0; i < li.length; i++) {
        const txtValue = (li[i].textContent || li[i].innerText).toUpperCase().trim();
        if (txtValue.indexOf(filter) > -1) {
            li[i].style.display = "";
        } else {
            li[i].style.display = "none";
        }
    }
}

function showModal(event) {
    const dialog = document.getElementById(event.getAttribute('id'));
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
window.filterKits = filterKits;
window.showModal = showModal;
window.filterMagazineLoadout = filterMagazineLoadout;
window.selectKit = selectKit;
window.closeToast = function () {
    if (previousToast) {
        clearTimeout(previousToast)
    }
    if (previousToastElement) {
        previousToastElement.classList.add("hidden")
    }
}

window.runtime.EventsOn('error', (e) => {
    document.getElementById('main').innerHTML = e;
})

window.runtime.EventsOn('toast.info', (e) => {
    showToast("success-toast", e, 5000);
})

window.runtime.EventsOn('toast.error', (e) => {
    showToast("error-toast", e, 10000);
})
window.runtime.EventsOn('clean_profile', (_e) => {
    document.getElementById('profile-selected').innerText = '';
})

function showToast(id, message, timeout = 2000) {
    if (previousToast) {
        clearTimeout(previousToast)
    }
    const toastElement = document.getElementById(id);
    previousToastElement = toastElement;
    const toastBody = toastElement.children.item(0).children.item(0);
    toastBody.innerText = message;
    toastElement.classList.remove("hidden")
    previousToast = setTimeout(() => {
        toastElement.classList.add("hidden");
        previousToastElement = null;
    }, timeout)
}

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
};
window.winterEvent();

