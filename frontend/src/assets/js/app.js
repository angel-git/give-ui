function setUsernameOnFooter() {
    const username = JSON.parse(document.getElementById("profile-selected-username").textContent);
    const nickname = JSON.parse(document.getElementById("profile-selected-nickname").textContent);
    document.getElementById('profile-selected').innerText = ": " + username + " - " + nickname;
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
        for (j = 0; j < innerLis.length; j++) {
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