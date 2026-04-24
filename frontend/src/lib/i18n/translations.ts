// Translation dictionary
import type { Locale } from "./locales";

export interface Translations {
    // App
    appName: string;
    appTitle: string;
    appCreator: string;

    // Navigation
    audioDevices: string;
    audio: string;
    themeEditor: string;
    settings: string;
    info: string;
    infoPanelTitle: string;
    infoPanelSubtitle: string;
    infoAboutAppTitle: string;
    infoReadmeTagline: string;
    infoAboutAppText: string;
    infoWhatItDoesTitle: string;
    infoWhatItDoesText: string;
    infoHowItWorksTitle: string;
    infoHowItWorksStep1: string;
    infoHowItWorksStep2: string;
    infoHowItWorksStep3: string;
    infoHowItWorksStep4: string;
    infoUsageTitle: string;
    infoUsageText: string;
    infoAnticheatTitle: string;
    infoAnticheatLead: string;
    infoAnticheatBullet1: string;
    infoAnticheatBullet2: string;
    infoAnticheatBullet3: string;
    infoAnticheatFooter: string;
    infoAppVersionLabel: string;
    infoDeveloperTitle: string;
    infoDeveloperLabel: string;
    infoDeveloperText: string;
    infoDonationsTitle: string;
    infoDonationsText: string;
    infoDonationsRussiaTitle: string;
    infoDonationsInternationalTitle: string;
    infoDonationsCryptoIntro: string;
    infoDonationCryptobotAction: string;
    infoDonationDonateStream: string;
    infoDonationCloudtips: string;
    infoDonationsFallback: string;
    infoWalletUsdtTrc: string;
    infoWalletUsdcPoly: string;
    infoWalletEth: string;
    infoWalletBtc: string;
    infoWalletLtc: string;
    contactsTitle: string;
    openTrustWallet: string;
    copied: string;
    copyAddress: string;

    // Audio Devices
    audioDevicesTitle: string;
    playbackDevice: string;
    defaultDevice: string;
    deviceInfo: string;
    pcm: string;
    rate: string;
    channels: string;
    defaultPeriod: string;
    minimumPeriod: string;
    latency: string;
    refreshDevices: string;
    collapseInfo: string;
    expandInfo: string;
    noDevicesFound: string;
    loadingDevices: string;
    errorLoadingDevices: string;
    deviceSelected: string;

    // Audio Player
    testSound: string;
    stopSound: string;
    playingSound: string;

    // Theme Editor
    themeEditorTitle: string;
    selectTheme: string;
    createNewTheme: string;
    themeName: string;
    enterThemeName: string;
    saveTheme: string;
    deleteTheme: string;
    applyTheme: string;
    themeSaved: string;
    themeDeleted: string;
    themeApplied: string;
    newTheme: string;
    cannotDeleteDefaultTheme: string;
    confirmDeleteTheme: string;
    geometrySection: string;
    sectors: string;
    rings: string;
    borderWidth: string;
    effectsSection: string;
    borderOpacity: string;
    sectionBaseOpacity: string;
    sectionBrightOpacity: string;
    sectionTimeout: string;
    blipSize: string;
    blipOpacity: string;
    blipTimeout: string;
    showBlips: string;
    themeBasicSettings: string;
    intensityMultiplier: string;
    colorsSection: string;
    positionSizeSection: string;
    sizeLabel: string;
    posX: string;
    posY: string;
    refreshThemes: string;
    themeInfo: string;
    noThemeSelected: string;
    themeNameRequired: string;
    themeColorAlpha: string;
    themeColorPreview: string;
    themeParamTooltipAria: string;
    themeUnitMs: string;
    themeHintThemeName: string;
    themeHintIntensityMultiplier: string;
    themeHintBackgroundColor: string;
    themeHintGridColor: string;
    themeHintSize: string;
    themeHintPosX: string;
    themeHintPosY: string;
    themeHintSectors: string;
    themeHintRings: string;
    themeHintBorderWidth: string;
    themeHintBorderOpacity: string;
    themeHintSectionBaseOpacity: string;
    themeHintSectionBrightOpacity: string;
    themeHintSectionTimeout: string;
    themeHintShowBlips: string;
    themeHintBlipSize: string;
    themeHintBlipOpacity: string;
    themeHintBlipTimeout: string;

    // Theme Colors
    backgroundColor: string;
    gridColor: string;
    sweepLineColor: string;
    blipColor: string;
    blipFadeColor: string;
    centerColor: string;
    textColor: string;

    // Radar
    startRadar: string;
    stopRadar: string;
    radarMode: string;
    normalMode: string;

    // Common
    close: string;
    cancel: string;
    confirm: string;
    error: string;
    success: string;
    loading: string;
    save: string;
    delete: string;
    edit: string;
    create: string;
    reset: string;
    yes: string;
    no: string;
    ok: string;
}

const en: Translations = {
    // App
    appName: "Sound Radar",
    appTitle: "See the sound, own the game. A real-time directional sound radar overlay for Windows",
    appCreator: "Created by Danil Solovyov",

    // Navigation
    audio: "Audio",
    audioDevices: "Audio Devices",
    themeEditor: "Themes",
    settings: "Settings",
    info: "Info",
    infoPanelTitle: "About the project",
    infoPanelSubtitle:
        "What the app does, how it works, anti-cheat notice, and ways to support the project.",
    infoAboutAppTitle: "Application",
    infoReadmeTagline:
        "A real-time directional sound visualizer (sound radar) and game overlay for Windows.",
    infoAboutAppText:
        "It helps gamers, including those with hearing loss, “see” 7.1 spatial audio on screen.",
    infoWhatItDoesTitle: "What it does",
    infoWhatItDoesText:
        "Transparent overlay on top of games, WASAPI loopback capture, up to 7.1 channel analysis, directional blip from channel vectors, customizable themes, low latency with adaptive polling.",
    infoHowItWorksTitle: "How it works",
    infoHowItWorksStep1:
        "Audio capture — WASAPI loopback reads peak levels for all channels of the selected device.",
    infoHowItWorksStep2:
        "Direction analysis — channel levels map to vectors (7.1 mapping); combined vector gives angle (atan2) and intensity.",
    infoHowItWorksStep3:
        "Smoothing — history buffer and intensity filter reduce noise.",
    infoHowItWorksStep4:
        "Visualization — Blip data is sent via Wails to the Svelte frontend and drawn on Canvas.",
    infoUsageTitle: "Usage",
    infoUsageText:
        "Pick a playback device — the radar starts automatically. Overlay modes: radar overlay (small window, click-through) and main overlay (fullscreen, same transparency). Themes: open Themes in the app. Tip: adjust radar position and size in the theme editor so it does not cover important game UI. Advanced users can edit config.toml next to the executable.",
    infoAnticheatTitle: "Anti-cheat notice",
    infoAnticheatLead:
        "Some anti-cheat systems (EAC, BattlEye, Vanguard) may falsely flag overlays as cheats.",
    infoAnticheatBullet1: "Does not interact with game memory.",
    infoAnticheatBullet2: "Does not inject DLLs.",
    infoAnticheatBullet3: "Does not emulate input.",
    infoAnticheatFooter:
        "The app only listens to audio (WASAPI) and draws a transparent window. Use at your own risk in online games. Game developers / anti-cheat vendors: contact me if you want to whitelist Game Radar (see contacts below).",
    infoAppVersionLabel: "Version",
    infoDeveloperTitle: "Developer",
    infoDeveloperLabel: "Development and support",
    infoDeveloperText: "Danil Solovyov.",
    infoDonationsTitle: "Support the project",
    infoDonationsText:
        "Game Radar is free and open source. If it helps you, you can support development.",
    infoDonationsRussiaTitle: "For users in Russia",
    infoDonationsInternationalTitle: "For users in other countries",
    infoDonationsCryptoIntro:
        "Cryptocurrency (no registration, direct transfer):",
    infoDonationCryptobotAction: "Send via CryptoBot",
    infoDonationDonateStream: "Donate.Stream (cards, SBP)",
    infoDonationCloudtips: "CloudTips (card transfer)",
    infoDonationsFallback:
        "If one method is unavailable in your region, use another from the list above (for example, direct USDT TRC20). Donations go toward maintenance, bug fixes, and new features.",
    infoWalletUsdtTrc: "USDT (recommended) · TRC20",
    infoWalletUsdcPoly: "USDC · Polygon",
    infoWalletEth: "ETH · Ethereum",
    infoWalletBtc: "Bitcoin · BTC",
    infoWalletLtc: "Litecoin · LTC",
    contactsTitle: "Contacts",
    openTrustWallet: "Open Trust Wallet",
    copied: "Copied",
    copyAddress: "Copy address",

    // Audio Devices
    audioDevicesTitle: "Audio Devices",
    playbackDevice: "Playback Device",
    defaultDevice: "Default",
    deviceInfo: "Device Information",
    pcm: "PCM",
    rate: "Rate",
    channels: "Channels",
    defaultPeriod: "Default Period",
    minimumPeriod: "Minimum Period",
    latency: "Latency",
    refreshDevices: "Refresh Devices",
    collapseInfo: "Hide Details",
    expandInfo: "Show Details",
    noDevicesFound: "No audio devices found",
    loadingDevices: "Loading devices...",
    errorLoadingDevices: "Failed to load devices",
    deviceSelected: "Device selected",

    // Audio Player
    testSound: "Test Sound",
    stopSound: "Stop",
    playingSound: "Playing...",

    // Theme Editor
    themeEditorTitle: "Theme Editor",
    selectTheme: "Select Theme",
    createNewTheme: "Create New Theme",
    themeName: "Theme Name",
    enterThemeName: "Enter theme name",
    saveTheme: "Save Theme",
    deleteTheme: "Delete Theme",
    applyTheme: "Apply Theme",
    themeSaved: "Theme saved successfully",
    themeDeleted: "Theme deleted",
    themeApplied: "Theme applied",
    newTheme: "New Theme",
    cannotDeleteDefaultTheme: "Cannot delete default theme",
    confirmDeleteTheme: "Are you sure you want to delete this theme?",
    geometrySection: "Geometry",
    sectors: "Sectors",
    rings: "Rings",
    borderWidth: "Border Width",
    effectsSection: "Effects",
    borderOpacity: "Border Opacity",
    sectionBaseOpacity: "Section Base Opacity",
    sectionBrightOpacity: "Transparency of sector illumination",
    sectionTimeout: "Section Timeout (ms)",
    blipSize: "Blip Size",
    blipOpacity: "Blip Opacity",
    blipTimeout: "Blip Timeout (ms)",
    showBlips: "Show Blips",
    themeBasicSettings: "Basic Settings",
    intensityMultiplier: "Sound gain",
    colorsSection: "Colors",
    positionSizeSection: "Position & Size",
    sizeLabel: "Size",
    posX: "X",
    posY: "Y",
    refreshThemes: "Refresh Themes",
    themeInfo: "Theme Info",
    noThemeSelected: "No theme selected",
    themeNameRequired: "Theme name cannot be empty",
    themeColorAlpha: "Opacity (A)",
    themeColorPreview: "Color preview",
    themeParamTooltipAria: "More about this setting",
    themeUnitMs: "ms",
    themeHintThemeName:
        "Name of this theme preset as it appears in the list and when saved. Use a short unique name so you can tell presets apart.",
    themeHintIntensityMultiplier:
        "Multiplies the strength of directional blips from the audio engine before they are drawn. Range 0.1–5. Lower values make quiet sounds subtler; higher values make the radar react more strongly. Does not change real audio volume.",
    themeHintBackgroundColor:
        "RGB from the color picker plus opacity (A, 0–255) for the canvas background behind the radar. Lower opacity makes the overlay more transparent so you can see the game through it.",
    themeHintGridColor:
        "Color and opacity of the radar grid (sectors and rings). A is 0–255; tune it so the grid is readable without cluttering the screen.",
    themeHintSize:
        "Radar diameter in pixels (120–1200). Larger is easier to read but covers more of the screen; match your resolution and HUD layout.",
    themeHintPosX:
        "Horizontal offset of the radar from the left edge of the overlay window in pixels. Adjust so the radar sits in a free area of your UI.",
    themeHintPosY:
        "Vertical offset from the top edge in pixels. Together with X and size, this places the radar where you want it.",
    themeHintSectors:
        "Number of angular sectors (direction bins). More sectors give finer direction resolution; fewer are simpler visually. Affects how sound direction is grouped on the wheel.",
    themeHintRings:
        "Number of concentric distance rings. More rings show distance in finer steps; fewer simplify the display.",
    themeHintBorderWidth:
        "Stroke width in pixels for sector and ring outlines. Thicker lines are easier to see; thinner lines look lighter.",
    themeHintBorderOpacity:
        "Opacity of grid outlines (0–1). Lower fades the border; 1 is fully opaque.",
    themeHintSectionBaseOpacity:
        "Resting opacity of sector fill (0–1). This is how visible inactive sectors are before a sound highlights one.",
    themeHintSectionBrightOpacity:
        "Peak opacity when a sector is highlighted by sound (0–1). Contrast with base opacity controls how obvious the flash is.",
    themeHintSectionTimeout:
        "How long a sector stays bright after activity, in milliseconds (50–5000). Shorter feels snappier; longer leaves a trail.",
    themeHintShowBlips:
        "When enabled, draws the moving directional blip dot in addition to sector highlights. Turn off if you only want sector glow.",
    themeHintBlipSize:
        "Diameter of the blip dot in pixels (1–40). Larger is easier to spot; smaller stays minimal on screen.",
    themeHintBlipOpacity:
        "Opacity of the blip graphic (0–1). Lower makes the dot softer; 1 is solid.",
    themeHintBlipTimeout:
        "How long the blip remains visible after the last update, in ms (50–5000). Controls fade persistence.",

    // Theme Colors
    backgroundColor: "Background Color",
    gridColor: "Grid Color",
    sweepLineColor: "Sweep Line Color",
    blipColor: "Blip Color",
    blipFadeColor: "Blip Fade Color",
    centerColor: "Center Color",
    textColor: "Text Color",

    // Radar
    startRadar: "Start Radar",
    stopRadar: "Stop Radar",
    radarMode: "Radar Mode",
    normalMode: "Normal Mode",

    // Common
    close: "Close",
    cancel: "Cancel",
    confirm: "Confirm",
    error: "Error",
    success: "Success",
    loading: "Loading...",
    save: "Save",
    delete: "Delete",
    edit: "Edit",
    create: "Create",
    reset: "Reset",
    yes: "Yes",
    no: "No",
    ok: "OK",
};

const ru: Translations = {
    // App
    appName: "Sound Radar",
    appTitle: "Видеть звук — владеть игрой. Радар звука в реальном времени для Windows",
    appCreator: "Создано Данилом Соловьёвым",

    // Navigation
    audio: "Аудио",
    audioDevices: "Аудио устройства",
    themeEditor: "Темы",
    settings: "Настройки",
    info: "Инфо",
    infoPanelTitle: "О проекте",
    infoPanelSubtitle:
        "Что делает приложение, как оно работает, предупреждение об античитах и способы поддержки.",
    infoAboutAppTitle: "Приложение",
    infoReadmeTagline:
        "Реалтайм-визуализатор направленного звука (звуковой радар) и игровой оверлей для Windows.",
    infoAboutAppText:
        "Помогает геймерам, особенно слабослышащим, «видеть» пространственное звучание 7.1 прямо на экране.",
    infoWhatItDoesTitle: "Возможности",
    infoWhatItDoesText:
        "Прозрачный оверлей поверх игр, захват через WASAPI loopback, анализ до 7.1, радар направления и силы звука по векторной сумме каналов, настраиваемые темы, низкая задержка и адаптивный опрос с учётом латентности устройства.",
    infoHowItWorksTitle: "Как это работает",
    infoHowItWorksStep1:
        "Захват звука — через WASAPI loopback считываются пиковые уровни всех каналов выбранного устройства.",
    infoHowItWorksStep2:
        "Анализ направления — уровни каналов превращаются в векторы (маппинг 7.1), суммарный вектор даёт угол (atan2) и интенсивность.",
    infoHowItWorksStep3:
        "Сглаживание — буфер истории и фильтр интенсивности убирают шумы.",
    infoHowItWorksStep4:
        "Визуализация — данные Blip передаются через Wails во фронтенд на Svelte, радар рисуется на Canvas.",
    infoUsageTitle: "Использование",
    infoUsageText:
        "Выберите устройство воспроизведения — радар запустится автоматически. Режимы оверлея: оверлей радара (небольшое окно, клики проходят сквозь) и главный оверлей (полноэкранно, та же прозрачность). Темы: вкладка «Темы» в приложении. Совет: настройте позицию и размер радара в редакторе тем, чтобы не перекрывать важный UI игры. Опытные пользователи могут править config.toml рядом с исполняемым файлом.",
    infoAnticheatTitle: "Предупреждение об античитах",
    infoAnticheatLead:
        "Некоторые античит-системы (EAC, BattlEye, Vanguard) могут ложно определять оверлей как читерскую программу.",
    infoAnticheatBullet1: "Не взаимодействует с памятью игры.",
    infoAnticheatBullet2: "Не внедряет DLL.",
    infoAnticheatBullet3: "Не эмулирует ввод.",
    infoAnticheatFooter:
        "Приложение только слушает звук (WASAPI) и рисует прозрачное окно. Используйте на свой страх и риск в онлайн-играх. Разработчикам игр и античит-системам: для внесения в белый список свяжитесь со мной (см. контакты ниже).",
    infoAppVersionLabel: "Версия",
    infoDeveloperTitle: "Разработчик",
    infoDeveloperLabel: "Разработка и поддержка",
    infoDeveloperText: "Данил Соловьёв.",
    infoDonationsTitle: "Поддержать проект",
    infoDonationsText:
        "Game Radar полностью бесплатен и открыт. Если он вам полезен, вы можете поддержать развитие проекта.",
    infoDonationsRussiaTitle: "Для пользователей из России",
    infoDonationsInternationalTitle: "Для пользователей из других стран",
    infoDonationsCryptoIntro:
        "Криптовалюта (без регистрации и без посредников):",
    infoDonationCryptobotAction: "Отправить через CryptoBot",
    infoDonationDonateStream: "Donate.Stream (карты, СБП)",
    infoDonationCloudtips: "CloudTips (перевод на карту)",
    infoDonationsFallback:
        "Если один способ недоступен в вашем регионе, используйте альтернативу из списка выше (например, прямой перевод в USDT TRC20). Пожертвования идут на поддержку проекта, исправления багов и новые функции.",
    infoWalletUsdtTrc: "USDT (рекомендуется) · TRC20",
    infoWalletUsdcPoly: "USDC · Polygon",
    infoWalletEth: "ETH · Ethereum",
    infoWalletBtc: "Bitcoin · BTC",
    infoWalletLtc: "Litecoin · LTC",
    contactsTitle: "Контакты",
    openTrustWallet: "Открыть в Trust Wallet",
    copied: "Скопировано",
    copyAddress: "Скопировать адрес",

    // Audio Devices
    audioDevicesTitle: "Аудио устройства",
    playbackDevice: "Устройство воспроизведения",
    defaultDevice: "По умолчанию",
    deviceInfo: "Информация об устройстве",
    pcm: "PCM",
    rate: "Частота",
    channels: "Каналы",
    defaultPeriod: "Период по умолчанию",
    minimumPeriod: "Минимальный период",
    latency: "Задержка",
    refreshDevices: "Обновить устройства",
    collapseInfo: "Скрыть детали",
    expandInfo: "Показать детали",
    noDevicesFound: "Аудио устройства не найдены",
    loadingDevices: "Загрузка устройств...",
    errorLoadingDevices: "Не удалось загрузить устройства",
    deviceSelected: "Устройство выбрано",

    // Audio Player
    testSound: "Тест звука",
    stopSound: "Стоп",
    playingSound: "Воспроизведение...",

    // Theme Editor
    themeEditorTitle: "Редактор темы",
    selectTheme: "Выбрать тему",
    createNewTheme: "Создать новую тему",
    themeName: "Название темы",
    enterThemeName: "Введите название темы",
    saveTheme: "Сохранить тему",
    deleteTheme: "Удалить тему",
    applyTheme: "Применить тему",
    themeSaved: "Тема успешно сохранена",
    themeDeleted: "Тема удалена",
    themeApplied: "Тема применена",
    newTheme: "Новая тема",
    cannotDeleteDefaultTheme: "Невозможно удалить тему по умолчанию",
    confirmDeleteTheme: "Вы уверены, что хотите удалить эту тему?",
    geometrySection: "Геометрия",
    sectors: "Секторы",
    rings: "Кольца",
    borderWidth: "Толщина границы",
    effectsSection: "Эффекты",
    borderOpacity: "Прозрачность границы",
    sectionBaseOpacity: "Базовая прозрачность сектора",
    sectionBrightOpacity: "Прозрачность подсветки сектора",
    sectionTimeout: "Таймаут сектора (мс)",
    blipSize: "Размер метки",
    blipOpacity: "Прозрачность метки",
    blipTimeout: "Таймаут метки (мс)",
    showBlips: "Показывать метки",
    themeBasicSettings: "Базовые настройки",
    intensityMultiplier: "Усиление звука",
    colorsSection: "Цвета",
    positionSizeSection: "Позиция и размер",
    sizeLabel: "Размер",
    posX: "X",
    posY: "Y",
    refreshThemes: "Обновить темы",
    themeInfo: "Информация о теме",
    noThemeSelected: "Тема не выбрана",
    themeNameRequired: "Название темы не может быть пустым",
    themeColorAlpha: "Прозрачность (A)",
    themeColorPreview: "Предпросмотр цвета",
    themeParamTooltipAria: "Подробнее об этой настройке",
    themeUnitMs: "мс",
    themeHintThemeName:
        "Внутреннее имя сохранённого пресета темы — так оно отображается в списке и при сохранении. Короткое уникальное имя упрощает выбор нужной темы.",
    themeHintIntensityMultiplier:
        "Множитель силы направленных меток (blip) после анализа звука, до отрисовки. Диапазон 0,1–5. Меньше — тише реакция на слабые звуки; больше — сильнее подсветка направления. Не меняет громкость реального аудио.",
    themeHintBackgroundColor:
        "RGB из палитры плюс непрозрачность канала A (0–255) для фона canvas под радаром. Меньше A — прозрачнее оверлей, лучше видно игру.",
    themeHintGridColor:
        "Цвет и прозрачность сетки радара (секторы и кольца). A от 0 до 255; подберите баланс читаемости и загромождения экрана.",
    themeHintSize:
        "Диаметр радара в пикселях (120–1200). Крупнее — удобнее читать, но больше закрывает экран; ориентируйтесь на разрешение и расположение HUD.",
    themeHintPosX:
        "Смещение радара по горизонтали от левого края окна оверлея в пикселях. Подберите свободную зону интерфейса.",
    themeHintPosY:
        "Смещение по вертикали от верхнего края в пикселях. Вместе с X и размером задаёт положение радара.",
    themeHintSectors:
        "Число угловых секторов (корзин направления). Больше — точнее направление, меньше — проще картинка. Влияет на группировку угла на «диске».",
    themeHintRings:
        "Число концентрических колец расстояния. Больше — более детальная шкала дальности, меньше — проще сетка.",
    themeHintBorderWidth:
        "Толщина линий границ секторов и колец в пикселях. Толще — заметнее, тоньше — легче визуально.",
    themeHintBorderOpacity:
        "Прозрачность линий сетки (0–1). Ниже — линии бледнее; 1 — непрозрачные.",
    themeHintSectionBaseOpacity:
        "Базовая заливка сектора в покое (0–1). Насколько видны неактивные секторы до подсветки звуком.",
    themeHintSectionBrightOpacity:
        "Максимальная непрозрачность сектора при подсветке (0–1). Контраст с базовой задаёт яркость вспышки.",
    themeHintSectionTimeout:
        "Сколько миллисекунд сектор остаётся подсвеченным после активности (50–5000). Короче — резче; дольше — длиннее «хвост».",
    themeHintShowBlips:
        "Включает движущуюся точку направления (метку) в дополнение к подсветке секторов. Выключите, если нужна только заливка секторов.",
    themeHintBlipSize:
        "Диаметр точки-метки в пикселях (1–40). Крупнее — заметнее, мельче — минимально на экране.",
    themeHintBlipOpacity:
        "Прозрачность метки (0–1). Ниже — мягче; 1 — плотная точка.",
    themeHintBlipTimeout:
        "Сколько миллисекунд метка остаётся видимой после последнего обновления (50–5000). Управляет затуханием.",

    // Theme Colors
    backgroundColor: "Цвет фона",
    gridColor: "Цвет сетки",
    sweepLineColor: "Цвет линии вращения",
    blipColor: "Цвет метки",
    blipFadeColor: "Цвет затухания метки",
    centerColor: "Центровой цвет",
    textColor: "Цвет текста",

    // Radar
    startRadar: "Запустить радар",
    stopRadar: "Остановить радар",
    radarMode: "Режим радара",
    normalMode: "Обычный режим",

    // Common
    close: "Закрыть",
    cancel: "Отмена",
    confirm: "Подтвердить",
    error: "Ошибка",
    success: "Успех",
    loading: "Загрузка...",
    save: "Сохранить",
    delete: "Удалить",
    edit: "Редактировать",
    create: "Создать",
    reset: "Сбросить",
    yes: "Да",
    no: "Нет",
    ok: "ОК",
};

export const TRANSLATIONS: Record<Locale, Translations> = {
    en,
    ru,
};
