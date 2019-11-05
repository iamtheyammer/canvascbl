const inject = `
<li class="menu-item ic-app-header__menu-list-item ">
    <a id="global_nav_conversations_link" href="https://canvascbl.com/?utm_source=extension" class="ic-app-header__menu-list-link" target="_blank">
        <div class="menu-item-icon-container">
            <span aria-hidden="true"><img src="https://canvascbl.com/logo-light-128.png" alt="canvascbl-logo">
            </span>
        </div>
        <div class="menu-item__text">CanvasCBL</div>
    </a>
</li>
`;

const el = document.getElementById("menu");
if (el) {
  el.insertAdjacentHTML("beforeend", inject);
}
