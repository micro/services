renderFootnotes = function () {
    const removeEl = (el) => {
        if (!el) return;
        el.remove ? el.remove() : el.parentNode.removeChild(el);
    };

    const insertAfter = (target, sib) => {
        target.after ? target.after(sib) : (
            target.parentNode.insertBefore(sib, target.nextSibling)
        );
    };

    const insideOut = (el) => {
        var p = el.parentNode, x = el.innerHTML,
            c = document.createElement('div');  // a tmp container
        insertAfter(p, c);
        c.appendChild(el);
        el.innerHTML = '';
        el.appendChild(p);
        p.innerHTML = x;  // let the original parent have the content of its child
        insertAfter(c, c.firstElementChild);
        removeEl(c);
    };

    document.querySelectorAll('.footnotes > ol > li[id^="fn"], #refs > div[id^="ref-"]').forEach(function (fn) {
        a = document.querySelectorAll('a[href="#' + fn.id + '"]');
        if (a.length === 0) return;
        a.forEach(function (el) { el.removeAttribute('href') });
        a = a[0];
        side = document.createElement('div');
        side.className = 'side side-right';
        if (/^fn/.test(fn.id)) {
            side.innerHTML = fn.innerHTML;
            var number = a.innerText;   // footnote number
            side.firstElementChild.innerHTML = '<span class="bg-number">' + number +
                '</span> ' + side.firstElementChild.innerHTML;
            removeEl(side.querySelector('a[href^="#fnref"]'));  // remove backreference
            a.parentNode.tagName === 'SUP' && insideOut(a);
        } else {
            side.innerHTML = fn.outerHTML;
            a = a.parentNode;
        }
        insertAfter(a, side);
        a.classList.add('note-ref');
        removeEl(fn);
    })
    document.querySelectorAll('.footnotes, #refs').forEach(function (fn) {
        var items = fn.children;
        if (fn.id === 'refs') return items.length === 0 && removeEl(fn);
        // there must be a <hr> and an <ol> left
        if (items.length !== 2 || items[0].tagName !== 'HR' || items[1].tagName !== 'OL') return;
        items[1].childElementCount === 0 && removeEl(fn);
    });
}();

renderAnchor = function () {
    for (let num = 1; num <= 6; num++) {
        // search h1-h6
        const headers = document.querySelectorAll('.article-post>h' + num);
        for (let i = 0; i < headers.length; i++) {
            const header = headers[i];
            // add anchor before h1-h6
            header.innerHTML = `<a href="#${header.id}" class="anchor"><svg class="icon" aria-hidden="true" focusable="false" height="16" version="1.1" viewBox="0 0 16 16" width="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg></a>${header.innerHTML}`;
        }
    }
}();

switchDarkMode = function () {
    const rootElement = document.documentElement; // <html>
    const darkModeStorageKey = 'user-color-scheme'; // use as localStorage's key
    const rootElementDarkModeAttributeName = 'data-user-color-scheme';
    const darkModeTogglebuttonElement = document.getElementById('dark-mode-button');

    const setLS = (k, v) => {
        try {
            localStorage.setItem(k, v);
        } catch (e) { }
    }

    const removeLS = (k) => {
        try {
            localStorage.removeItem(k);
        } catch (e) { }
    }

    const getLS = (k) => {
        try {
            return localStorage.getItem(k);
        } catch (e) {
            return null // the same as localStorage.getItem() get nothing
        }
    }

    const getModeFromCSSMediaQuery = () => {
        // use matchMedia API
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }

    const resetRootDarkModeAttributeAndLS = () => {
        rootElement.removeAttribute(rootElementDarkModeAttributeName);
        removeLS(darkModeStorageKey);
    }

    const validColorModeKeys = {
        'dark': true,
        'light': true
    }

    const modeIcons = {
        'dark': '☀️',
        'light': '🌙'
    }

    const setModeButtonIcon = (mode) => {
        darkModeTogglebuttonElement.innerHTML = modeIcons[mode]
    }

    const applyCustomDarkModeSettings = (mode) => {
        // receive user's operation or get previous mode from localStorage
        const currentSetting = mode || getLS(darkModeStorageKey);

        if (currentSetting === getModeFromCSSMediaQuery()) {
            // When the user selected mode equal prefers-color-scheme 
            // reset and restored to automatic mode
            nowMode = getModeFromCSSMediaQuery()
            resetRootDarkModeAttributeAndLS();
        } else if (validColorModeKeys[currentSetting]) {
            nowMode = currentSetting
            rootElement.setAttribute(rootElementDarkModeAttributeName, currentSetting);
        } else {
            // 首次访问或从未使用过开关、localStorage 中没有存储的值，currentSetting 是 null
            // 或者 localStorage 被篡改，currentSetting 不是合法值
            nowMode = getModeFromCSSMediaQuery()
            resetRootDarkModeAttributeAndLS();
        }
        setModeButtonIcon(nowMode)
    }

    const invertDarkModeObj = {
        'dark': 'light',
        'light': 'dark'
    }

    const toggleCustomDarkMode = () => {
        let currentSetting = getLS(darkModeStorageKey);

        if (validColorModeKeys[currentSetting]) {
            // get mode from localStorage and set the opposite
            currentSetting = invertDarkModeObj[currentSetting];
        } else if (currentSetting === null) {
            // if get null from localStorage
            // get mode from prefers-color-scheme and set the opposite
            currentSetting = invertDarkModeObj[getModeFromCSSMediaQuery()];
        } else {
            // get anything error, return
            return;
        }
        // set opposite mode into localStorage
        setLS(darkModeStorageKey, currentSetting);

        return currentSetting;
    }

    // when page loaded set page mode
    applyCustomDarkModeSettings();

    darkModeTogglebuttonElement.addEventListener('click', () => {
        // handle user click switch dark mode button
        applyCustomDarkModeSettings(toggleCustomDarkMode());
    })
}();