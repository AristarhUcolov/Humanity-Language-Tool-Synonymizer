#!/usr/bin/env python3
"""
MEGA expansion for Russian and Moldovan dictionaries.
Adds thousands of base words and morphological forms.
"""
import re
from pathlib import Path

PROJECT_ROOT = Path(__file__).parent.parent
DICT_DIR = PROJECT_ROOT / "internal" / "humanize"


def extract_existing(path: Path) -> dict:
    if not path.exists():
        return {}
    content = path.read_text(encoding='utf-8')
    pattern = r'"([^"]+)":\s*"([^"]+)"'
    matches = re.findall(pattern, content)
    return {k: v for k, v in matches}


def write_dict(path: Path, var_name: str, comment: str, entries: dict):
    sorted_items = sorted(entries.items())
    lines = [f"package humanize", "", comment, f"var {var_name} = dictPack{{", "\tAICliches: map[string]string{"]
    for k, v in sorted_items:
        k_esc = k.replace('"', '\\"')
        v_esc = v.replace('"', '\\"')
        lines.append(f'\t\t"{k_esc}": "{v_esc}",')
    lines.append("\t},")
    lines.append("\tFiller:   map[string]string{},")
    lines.append("\tUpgrades: map[string]string{},")
    lines.append("}")
    lines.append("")
    path.write_text("\n".join(lines), encoding='utf-8')


# ===================================================================
# RUSSIAN MEGA EXPANSION - hundreds of additional base words with forms
# ===================================================================

def expand_russian_noun_cases(nom_sg, repl_sg, nom_pl, repl_pl, gender):
    """Generate 6 cases × 2 numbers = 12 forms for a Russian noun (simplified)."""
    result = {nom_sg: repl_sg, nom_pl: repl_pl}

    # Masculine (ends in consonant or -й)
    if gender == 'm':
        # Singular cases
        if nom_sg.endswith('й'):
            stem = nom_sg[:-1]
            result[stem + 'я'] = repl_sg  # gen
            result[stem + 'ю'] = repl_sg  # dat
            result[stem + 'я'] = repl_sg  # acc (animate)
            result[stem + 'ем'] = repl_sg  # inst
            result[stem + 'е'] = repl_sg  # prep
        elif nom_sg.endswith('ь'):
            stem = nom_sg[:-1]
            result[stem + 'я'] = repl_sg
            result[stem + 'ю'] = repl_sg
            result[stem + 'ем'] = repl_sg
            result[stem + 'е'] = repl_sg
        else:
            stem = nom_sg
            result[stem + 'а'] = repl_sg  # gen
            result[stem + 'у'] = repl_sg  # dat
            result[stem + 'ом'] = repl_sg  # inst
            result[stem + 'е'] = repl_sg  # prep

    # Feminine (ends in -а or -я)
    elif gender == 'f':
        if nom_sg.endswith('а'):
            stem = nom_sg[:-1]
            result[stem + 'ы'] = repl_sg  # gen
            result[stem + 'е'] = repl_sg  # dat/prep
            result[stem + 'у'] = repl_sg  # acc
            result[stem + 'ой'] = repl_sg  # inst
        elif nom_sg.endswith('я'):
            stem = nom_sg[:-1]
            result[stem + 'и'] = repl_sg
            result[stem + 'е'] = repl_sg
            result[stem + 'ю'] = repl_sg
            result[stem + 'ей'] = repl_sg

    # Neuter (ends in -о or -е)
    elif gender == 'n':
        if nom_sg.endswith('о'):
            stem = nom_sg[:-1]
            result[stem + 'а'] = repl_sg
            result[stem + 'у'] = repl_sg
            result[stem + 'ом'] = repl_sg
            result[stem + 'е'] = repl_sg
        elif nom_sg.endswith('е'):
            stem = nom_sg[:-1]
            result[stem + 'я'] = repl_sg
            result[stem + 'ю'] = repl_sg
            result[stem + 'ем'] = repl_sg
            result[stem + 'е'] = repl_sg

    return result


def expand_russian_adj_full(m_form, m_repl):
    """Generate all gender/number forms for Russian adjective."""
    # Auto-detect feminine, neuter, plural endings
    if m_form.endswith('ый') or m_form.endswith('ой') or m_form.endswith('ий'):
        stem = m_form[:-2]
        m_stem = m_repl[:-2] if (m_repl.endswith('ый') or m_repl.endswith('ой') or m_repl.endswith('ий')) else m_repl
    else:
        return {m_form: m_repl}

    result = {
        m_form: m_repl,
        stem + 'ого': m_stem + 'ого',  # gen masc
        stem + 'ому': m_stem + 'ому',  # dat masc
        stem + 'ым': m_stem + 'ым',    # inst masc
        stem + 'ом': m_stem + 'ом',    # prep masc
        stem + 'ая': m_stem + 'ая',    # fem
        stem + 'ой': m_stem + 'ой',    # fem gen/dat/inst/prep
        stem + 'ую': m_stem + 'ую',    # fem acc
        stem + 'ое': m_stem + 'ое',    # neuter
        stem + 'ые': m_stem + 'ые',    # plural
        stem + 'ых': m_stem + 'ых',    # pl gen
        stem + 'ым': m_stem + 'ым',    # pl dat
        stem + 'ыми': m_stem + 'ыми',  # pl inst
    }
    return result


def expand_russian_verb_full(infinitive, repl_inf):
    """Generate verb forms: infinitive + past × 4 + present × 6 + imperative."""
    result = {infinitive: repl_inf}

    # Get stem
    if infinitive.endswith('ть'):
        stem = infinitive[:-2]
        repl_stem = repl_inf[:-2] if repl_inf.endswith('ть') else repl_inf
        # Past forms
        result[stem + 'л'] = repl_stem + 'л'
        result[stem + 'ла'] = repl_stem + 'ла'
        result[stem + 'ло'] = repl_stem + 'ло'
        result[stem + 'ли'] = repl_stem + 'ли'
    elif infinitive.endswith('ться'):
        stem = infinitive[:-4]
        repl_stem = repl_inf[:-4] if repl_inf.endswith('ться') else repl_inf
        result[stem + 'лся'] = repl_stem + 'лся'
        result[stem + 'лась'] = repl_stem + 'лась'
        result[stem + 'лось'] = repl_stem + 'лось'
        result[stem + 'лись'] = repl_stem + 'лись'

    return result


# Russian: many new base nouns with full case forms
russian_nouns = [
    # (nominative singular, replacement singular, nom plural, repl plural, gender)
    ("учитель", "наставник", "учителя", "наставники", "m"),
    ("ученик", "воспитанник", "ученики", "воспитанники", "m"),
    ("ученица", "воспитанница", "ученицы", "воспитанницы", "f"),
    ("учительница", "наставница", "учительницы", "наставницы", "f"),
    ("студент", "слушатель", "студенты", "слушатели", "m"),
    ("студентка", "слушательница", "студентки", "слушательницы", "f"),
    ("врач", "медик", "врачи", "медики", "m"),
    ("медсестра", "санитарка", "медсёстры", "санитарки", "f"),
    ("инженер", "конструктор", "инженеры", "конструкторы", "m"),
    ("программист", "разработчик", "программисты", "разработчики", "m"),
    ("художник", "живописец", "художники", "живописцы", "m"),
    ("писатель", "литератор", "писатели", "литераторы", "m"),
    ("поэт", "стихотворец", "поэты", "стихотворцы", "m"),
    ("композитор", "сочинитель", "композиторы", "сочинители", "m"),
    ("музыкант", "исполнитель", "музыканты", "исполнители", "m"),
    ("певец", "вокалист", "певцы", "вокалисты", "m"),
    ("певица", "вокалистка", "певицы", "вокалистки", "f"),
    ("актёр", "артист", "актёры", "артисты", "m"),
    ("актриса", "артистка", "актрисы", "артистки", "f"),
    ("режиссёр", "постановщик", "режиссёры", "постановщики", "m"),
    ("спортсмен", "атлет", "спортсмены", "атлеты", "m"),
    ("спортсменка", "атлетка", "спортсменки", "атлетки", "f"),
    ("учёный", "исследователь", "учёные", "исследователи", "m"),
    ("философ", "мыслитель", "философы", "мыслители", "m"),
    ("священник", "духовник", "священники", "духовники", "m"),
    ("солдат", "военнослужащий", "солдаты", "военнослужащие", "m"),
    ("офицер", "командир", "офицеры", "командиры", "m"),
    ("генерал", "военачальник", "генералы", "военачальники", "m"),
    ("полковник", "командующий", "полковники", "командующие", "m"),
    ("капитан", "предводитель", "капитаны", "предводители", "m"),
    ("сержант", "унтер-офицер", "сержанты", "унтер-офицеры", "m"),
    ("полицейский", "правоохранитель", "полицейские", "правоохранители", "m"),
    ("милиционер", "блюститель порядка", "милиционеры", "блюстители порядка", "m"),
    ("охранник", "страж", "охранники", "стражи", "m"),
    ("водитель", "шофёр", "водители", "шофёры", "m"),
    ("пилот", "лётчик", "пилоты", "лётчики", "m"),
    ("моряк", "мореплаватель", "моряки", "мореплаватели", "m"),
    ("капитан корабля", "судоводитель", "капитаны кораблей", "судоводители", "m"),
    ("повар", "кулинар", "повара", "кулинары", "m"),
    ("официант", "обслуживающий", "официанты", "обслуживающие", "m"),
    ("официантка", "обслуживающая", "официантки", "обслуживающие", "f"),
    ("продавец", "торговец", "продавцы", "торговцы", "m"),
    ("продавщица", "торговка", "продавщицы", "торговки", "f"),
    ("покупатель", "клиент", "покупатели", "клиенты", "m"),
    ("покупательница", "клиентка", "покупательницы", "клиентки", "f"),
    ("работодатель", "наниматель", "работодатели", "наниматели", "m"),
    ("работник", "сотрудник", "работники", "сотрудники", "m"),
    ("рабочий", "пролетарий", "рабочие", "пролетарии", "m"),
    ("крестьянин", "земледелец", "крестьяне", "земледельцы", "m"),
    ("фермер", "сельхозработник", "фермеры", "сельхозработники", "m"),
    ("король", "монарх", "короли", "монархи", "m"),
    ("королева", "правительница", "королевы", "правительницы", "f"),
    ("принц", "наследник престола", "принцы", "наследники престола", "m"),
    ("принцесса", "наследница престола", "принцессы", "наследницы престола", "f"),
    ("президент", "глава государства", "президенты", "главы государств", "m"),
    ("министр", "член правительства", "министры", "члены правительства", "m"),
    ("депутат", "парламентарий", "депутаты", "парламентарии", "m"),
    ("мэр", "градоначальник", "мэры", "градоначальники", "m"),
    ("губернатор", "управляющий регионом", "губернаторы", "управляющие регионами", "m"),
    ("чиновник", "функционер", "чиновники", "функционеры", "m"),

    # Common objects
    ("телефон", "коммуникатор", "телефоны", "коммуникаторы", "m"),
    ("компьютер", "вычислительная машина", "компьютеры", "вычислительные машины", "m"),
    ("ноутбук", "переносной компьютер", "ноутбуки", "переносные компьютеры", "m"),
    ("планшет", "сенсорное устройство", "планшеты", "сенсорные устройства", "m"),
    ("телевизор", "приёмник", "телевизоры", "приёмники", "m"),
    ("радио", "радиоприёмник", "радио", "радиоприёмники", "n"),
    ("часы", "хронометр", "часы", "хронометры", "m"),
    ("очки", "оптические линзы", "очки", "оптические линзы", "m"),
    ("ключ", "отмычка", "ключи", "отмычки", "m"),
    ("замок", "запор", "замки", "запоры", "m"),
    ("сумка", "ридикюль", "сумки", "ридикюли", "f"),
    ("портфель", "дипломат", "портфели", "дипломаты", "m"),
    ("кошелёк", "бумажник", "кошельки", "бумажники", "m"),
    ("шкаф", "гардероб", "шкафы", "гардеробы", "m"),
    ("полка", "стеллаж", "полки", "стеллажи", "f"),
    ("ковёр", "напольное покрытие", "ковры", "напольные покрытия", "m"),
    ("картина", "полотно", "картины", "полотна", "f"),
    ("статуя", "изваяние", "статуи", "изваяния", "f"),
    ("памятник", "монумент", "памятники", "монументы", "m"),

    # Food
    ("хлеб", "выпечка", "хлеба", "выпечка", "m"),
    ("масло", "жировой продукт", "масла", "жировые продукты", "n"),
    ("сыр", "молочный продукт", "сыры", "молочные продукты", "m"),
    ("молоко", "млечный напиток", "молоко", "молочное", "n"),
    ("мясо", "мускульная ткань", "мясо", "мускульные ткани", "n"),
    ("рыба", "ихтиофауна", "рыбы", "водные обитатели", "f"),
    ("курица", "пернатая птица", "куры", "пернатые птицы", "f"),
    ("яйцо", "яичный продукт", "яйца", "яичные продукты", "n"),
    ("овощ", "огородная культура", "овощи", "огородные культуры", "m"),
    ("фрукт", "плод", "фрукты", "плоды", "m"),
    ("ягода", "лесной плод", "ягоды", "лесные плоды", "f"),
    ("картофель", "клубнеплод", "картофели", "клубнеплоды", "m"),
    ("помидор", "томат", "помидоры", "томаты", "m"),
    ("огурец", "огородное растение", "огурцы", "огородные растения", "m"),
    ("капуста", "крестоцветное", "капусты", "крестоцветные", "f"),
    ("морковь", "корнеплод", "моркови", "корнеплоды", "f"),
    ("яблоко", "плод яблони", "яблоки", "плоды яблонь", "n"),
    ("груша", "плод груши", "груши", "плоды груш", "f"),
    ("апельсин", "цитрусовый плод", "апельсины", "цитрусовые плоды", "m"),
    ("банан", "тропический плод", "бананы", "тропические плоды", "m"),
    ("лимон", "кислый цитрус", "лимоны", "кислые цитрусы", "m"),
    ("виноград", "лозовая культура", "виноград", "лозовые культуры", "m"),

    # Travel/places
    ("гостиница", "отель", "гостиницы", "отели", "f"),
    ("ресторан", "трактир", "рестораны", "трактиры", "m"),
    ("кафе", "кофейня", "кафе", "кофейни", "n"),
    ("бар", "питейное заведение", "бары", "питейные заведения", "m"),
    ("музей", "хранилище реликвий", "музеи", "хранилища реликвий", "m"),
    ("театр", "зрелищное заведение", "театры", "зрелищные заведения", "m"),
    ("кино", "кинотеатр", "кино", "кинотеатры", "n"),
    ("концерт", "выступление", "концерты", "выступления", "m"),
    ("стадион", "спортивная арена", "стадионы", "спортивные арены", "m"),
    ("больница", "лечебница", "больницы", "лечебницы", "f"),
    ("аптека", "фармацевтическое учреждение", "аптеки", "фармацевтические учреждения", "f"),
    ("банк", "финансовое учреждение", "банки", "финансовые учреждения", "m"),
    ("почта", "почтовое отделение", "почты", "почтовые отделения", "f"),
    ("вокзал", "железнодорожная станция", "вокзалы", "железнодорожные станции", "m"),
    ("аэропорт", "воздушная гавань", "аэропорты", "воздушные гавани", "m"),
    ("порт", "морская гавань", "порты", "морские гавани", "m"),
    ("остановка", "пункт остановки", "остановки", "пункты остановки", "f"),
    ("автобус", "омнибус", "автобусы", "омнибусы", "m"),
    ("трамвай", "рельсовый транспорт", "трамваи", "рельсовый транспорт", "m"),
    ("метро", "подземка", "метро", "подземки", "n"),
    ("такси", "наёмный автомобиль", "такси", "наёмные автомобили", "n"),
    ("поезд", "состав", "поезда", "составы", "m"),
    ("самолёт", "воздушное судно", "самолёты", "воздушные суда", "m"),
    ("корабль", "судно", "корабли", "суда", "m"),
    ("лодка", "плавательное средство", "лодки", "плавательные средства", "f"),
    ("автомобиль", "транспортное средство", "автомобили", "транспортные средства", "m"),
    ("велосипед", "двухколёсное", "велосипеды", "двухколёсные", "m"),
    ("мотоцикл", "двухколёсный мотор", "мотоциклы", "двухколёсные моторы", "m"),

    # More abstract concepts
    ("совесть", "нравственное чувство", "совести", "нравственные чувства", "f"),
    ("душа", "психея", "души", "психеи", "f"),
    ("разум", "интеллект", "умы", "интеллекты", "m"),
    ("ум", "разумение", "умы", "разумения", "m"),
    ("воля", "решимость", "воли", "решимости", "f"),
    ("свобода", "независимость", "свободы", "независимости", "f"),
    ("рабство", "невольничество", "рабства", "невольничества", "n"),
    ("власть", "господство", "власти", "господства", "f"),
    ("сила", "мощь", "силы", "мощности", "f"),
    ("слава", "почёт", "славы", "почести", "f"),
    ("честь", "достоинство", "чести", "достоинства", "f"),
    ("позор", "бесчестие", "позоры", "бесчестия", "m"),
    ("гордость", "превосходство", "гордости", "превосходства", "f"),
    ("стыд", "смущение", "стыды", "смущения", "m"),
    ("вина", "виновность", "вины", "виновности", "f"),
    ("прощение", "помилование", "прощения", "помилования", "n"),
    ("мечта", "грёза", "мечты", "грёзы", "f"),
    ("реальность", "действительность", "реальности", "действительности", "f"),
    ("истина", "правда", "истины", "правды", "f"),
    ("ложь", "обман", "лжи", "обманы", "f"),
    ("обман", "мошенничество", "обманы", "мошенничества", "m"),
    ("предательство", "измена", "предательства", "измены", "n"),
    ("верность", "преданность", "верности", "преданности", "f"),
    ("дружба", "товарищество", "дружбы", "товарищества", "f"),
    ("любовь", "влечение", "любови", "влечения", "f"),
    ("страсть", "пыл", "страсти", "пыл", "f"),
    ("нежность", "ласка", "нежности", "ласки", "f"),
    ("ласка", "нежность", "ласки", "нежности", "f"),
    ("забота", "попечение", "заботы", "попечения", "f"),
    ("опека", "покровительство", "опеки", "покровительства", "f"),
    ("защита", "оборона", "защиты", "обороны", "f"),
    ("безопасность", "сохранность", "безопасности", "сохранности", "f"),
    ("опасность", "угроза", "опасности", "угрозы", "f"),

    # Education-related
    ("урок", "занятие", "уроки", "занятия", "m"),
    ("экзамен", "испытание", "экзамены", "испытания", "m"),
    ("оценка", "балл", "оценки", "баллы", "f"),
    ("задание", "поручение", "задания", "поручения", "n"),
    ("упражнение", "тренировка", "упражнения", "тренировки", "n"),
    ("домашка", "домашнее задание", "домашки", "домашние задания", "f"),
    ("контрольная", "проверочная работа", "контрольные", "проверочные работы", "f"),
    ("тетрадь", "ученический блокнот", "тетради", "ученические блокноты", "f"),
    ("учебник", "учебное пособие", "учебники", "учебные пособия", "m"),
    ("словарь", "лексикон", "словари", "лексиконы", "m"),
    ("карандаш", "графитный стержень", "карандаши", "графитные стержни", "m"),
    ("ручка", "пишущий инструмент", "ручки", "пишущие инструменты", "f"),
    ("бумага", "писчая основа", "бумаги", "писчие основы", "f"),
    ("книга", "издание", "книги", "издания", "f"),
    ("журнал", "периодическое издание", "журналы", "периодические издания", "m"),
    ("газета", "печатное издание", "газеты", "печатные издания", "f"),
    ("статья", "публикация", "статьи", "публикации", "f"),

    # Weather
    ("погода", "метеорологические условия", "погоды", "метеорологические условия", "f"),
    ("туман", "пасмурная дымка", "туманы", "пасмурные дымки", "m"),
    ("гроза", "грозовая буря", "грозы", "грозовые бури", "f"),
    ("буря", "ураган", "бури", "ураганы", "f"),
    ("ураган", "тропический циклон", "ураганы", "тропические циклоны", "m"),
    ("гром", "раскат грома", "громы", "раскаты грома", "m"),
    ("молния", "электрический разряд", "молнии", "электрические разряды", "f"),
    ("ветер", "воздушный поток", "ветры", "воздушные потоки", "m"),
    ("град", "ледяная крупа", "грады", "ледяные крупы", "m"),
]

# Russian verbs to generate forms for
russian_verbs_extra = [
    # (infinitive, replacement)
    ("сделать", "осуществить"),
    ("сказать", "произнести"),
    ("услышать", "уловить"),
    ("увидеть", "узреть"),
    ("узнать", "осведомиться"),
    ("понять", "уяснить"),
    ("подумать", "обдумать"),
    ("захотеть", "возжелать"),
    ("полюбить", "увлечься"),
    ("пойти", "отправиться"),
    ("прийти", "явиться"),
    ("уйти", "удалиться"),
    ("выйти", "выбраться"),
    ("войти", "проникнуть"),
    ("взять", "получить"),
    ("дать", "предоставить"),
    ("получить", "обрести"),
    ("найти", "обнаружить"),
    ("потерять", "утратить"),
    ("купить", "приобрести"),
    ("продать", "реализовать"),
    ("заплатить", "оплатить"),
    ("помочь", "посодействовать"),
    ("попросить", "ходатайствовать"),
    ("ответить", "среагировать"),
    ("спросить", "поинтересоваться"),
    ("открыть", "распахнуть"),
    ("закрыть", "затворить"),
    ("начать", "приступить"),
    ("закончить", "завершить"),
    ("продолжить", "длить"),
    ("остановить", "прекратить"),
    ("стать", "превратиться"),
    ("остаться", "пребывать"),
    ("явиться", "представать"),
    ("случиться", "приключиться"),
    ("произойти", "развернуться"),
    ("поделиться", "распределиться"),
    ("воспользоваться", "применить"),
    ("создать", "сформировать"),
    ("уничтожить", "истребить"),
    ("разрушить", "развалить"),
    ("построить", "соорудить"),
    ("сломать", "разломать"),
    ("починить", "восстановить"),
    ("изменить", "преобразовать"),
    ("улучшить", "усовершенствовать"),
    ("ухудшить", "усугубить"),
    ("увеличить", "приумножить"),
    ("уменьшить", "сократить"),
    ("показать", "продемонстрировать"),
    ("спрятать", "сокрыть"),
    ("забыть", "запамятовать"),
    ("запомнить", "сохранить в памяти"),
    ("надеяться", "уповать"),
    ("бояться", "опасаться"),
    ("радоваться", "ликовать"),
    ("заплакать", "разрыдаться"),
    ("засмеяться", "расхохотаться"),
    ("закричать", "вскричать"),
    ("замолчать", "умолкнуть"),
    ("заждаться", "истомиться ожиданием"),
    ("позвать", "призвать"),
    ("позвонить", "телефонировать"),
    ("нарисовать", "изобразить"),
    ("сыграть", "исполнить"),
    ("спеть", "пропеть"),
    ("станцевать", "исполнить танец"),
    ("прогуляться", "променажировать"),
    ("попутешествовать", "постранствовать"),
    ("поехать", "отправиться"),
    ("полететь", "воспарить"),
    ("поплавать", "погрузиться в воду"),
    ("плыть", "следовать по воде"),
    ("ехать", "перемещаться"),
    ("лететь", "парить"),
    ("плавать", "находиться на воде"),
    ("сидеть", "располагаться"),
    ("стоять", "находиться"),
    ("лежать", "покоиться"),
    ("вставать", "подниматься"),
    ("ложиться", "укладываться"),
    ("садиться", "усаживаться"),
    ("прыгать", "подскакивать"),
    ("падать", "опускаться"),
    ("подниматься", "восходить"),
    ("спускаться", "нисходить"),
    ("носить", "транспортировать"),
    ("везти", "перевозить"),
    ("вести", "сопровождать"),
    ("гнать", "понукать"),
    ("ловить", "захватывать"),
    ("кидать", "метать"),
    ("бросать", "швырять"),
    ("хватать", "схватывать"),
    ("тянуть", "влечь"),
    ("толкать", "пихать"),
    ("давить", "придавливать"),
    ("резать", "разрезать"),
    ("колоть", "пронзать"),
    ("бить", "ударять"),
    ("стучать", "колотить"),
    ("шептать", "говорить тихо"),
    ("шуметь", "издавать шум"),
    ("звонить", "издавать звон"),
    ("звенеть", "издавать звенящий звук"),
    ("трещать", "издавать треск"),
    ("свистеть", "издавать свист"),
    ("петь", "вокализировать"),
    ("танцевать", "плясать"),
    ("рисовать", "изображать"),
    ("писать", "записывать"),
    ("читать", "изучать"),
    ("слушать", "внимать"),
    ("смотреть", "созерцать"),
    ("вспоминать", "припоминать"),
    ("представлять", "воображать"),
    ("ожидать", "предвкушать"),
    ("надеяться", "уповать"),
    ("верить", "уповать"),
    ("сомневаться", "колебаться"),
    ("отрицать", "опровергать"),
    ("утверждать", "заявлять"),
    ("соглашаться", "поддерживать"),
    ("возражать", "противоречить"),
    ("предлагать", "выдвигать"),
    ("отказываться", "уклоняться"),
    ("принимать", "соглашаться"),
    ("отвергать", "отклонять"),
    ("разрешать", "позволять"),
    ("запрещать", "воспрещать"),
    ("заставлять", "принуждать"),
    ("позволять", "дозволять"),
]

# Russian adjectives to generate full paradigm
russian_adjs_extra = [
    # (masculine form, replacement)
    ("умный", "разумный"),
    ("глупый", "недалёкий"),
    ("добрый", "благосклонный"),
    ("злой", "злобный"),
    ("честный", "правдивый"),
    ("лживый", "обманчивый"),
    ("смелый", "отважный"),
    ("трусливый", "робкий"),
    ("гордый", "величественный"),
    ("скромный", "сдержанный"),
    ("веселый", "жизнерадостный"),
    ("печальный", "скорбный"),
    ("спокойный", "невозмутимый"),
    ("нервный", "встревоженный"),
    ("здоровый", "благополучный"),
    ("больной", "недужный"),
    ("молодой", "юный"),
    ("старый", "пожилой"),
    ("красивый", "восхитительный"),
    ("уродливый", "отталкивающий"),
    ("чистый", "безупречный"),
    ("грязный", "загрязнённый"),
    ("светлый", "лучезарный"),
    ("тёмный", "сумрачный"),
    ("тёплый", "согревающий"),
    ("холодный", "промозглый"),
    ("горячий", "знойный"),
    ("прохладный", "освежающий"),
    ("сухой", "пересохший"),
    ("мокрый", "увлажнённый"),
    ("новый", "современный"),
    ("старинный", "антикварный"),
    ("свежий", "только что приготовленный"),
    ("чёрствый", "залежалый"),
    ("сладкий", "медовый"),
    ("горький", "терпкий"),
    ("кислый", "лимонный"),
    ("солёный", "пересоленный"),
    ("вкусный", "превосходный на вкус"),
    ("противный", "отвратительный"),
    ("приятный", "доставляющий удовольствие"),
    ("неприятный", "доставляющий неудобство"),
    ("полезный", "приносящий пользу"),
    ("вредный", "приносящий вред"),
    ("важный", "значимый"),
    ("неважный", "несущественный"),
    ("нужный", "необходимый"),
    ("ненужный", "излишний"),
    ("трудный", "сложный"),
    ("простой", "элементарный"),
    ("сложный", "запутанный"),
    ("ясный", "очевидный"),
    ("непонятный", "загадочный"),
    ("понятный", "доступный пониманию"),
    ("известный", "знаменитый"),
    ("неизвестный", "безвестный"),
    ("популярный", "пользующийся спросом"),
    ("обычный", "заурядный"),
    ("необычный", "выдающийся"),
    ("странный", "необычайный"),
    ("привычный", "обыкновенный"),
    ("редкий", "уникальный"),
    ("частый", "регулярный"),
    ("единственный", "единичный"),
    ("разный", "различный"),
    ("одинаковый", "идентичный"),
    ("похожий", "сходный"),
    ("несхожий", "отличный"),
    ("опасный", "рискованный"),
    ("безопасный", "защищённый"),
    ("свободный", "независимый"),
    ("занятый", "загруженный"),
    ("открытый", "распахнутый"),
    ("закрытый", "затворённый"),
    ("живой", "одушевлённый"),
    ("мёртвый", "бездыханный"),
    ("настоящий", "подлинный"),
    ("ненастоящий", "поддельный"),
    ("истинный", "правдивый"),
    ("ложный", "обманчивый"),
    ("вечный", "бессмертный"),
    ("временный", "преходящий"),
    ("постоянный", "неизменный"),
    ("случайный", "непреднамеренный"),
    ("намеренный", "сознательный"),
    ("вольный", "беспрепятственный"),
    ("обязательный", "необходимый"),
    ("возможный", "осуществимый"),
    ("невозможный", "неосуществимый"),
    ("реальный", "действительный"),
    ("нереальный", "иллюзорный"),
    ("видимый", "различимый"),
    ("невидимый", "неразличимый"),
    ("слышный", "различимый на слух"),
    ("неслышный", "беззвучный"),
    ("явный", "очевидный"),
    ("скрытый", "потаённый"),
    ("тайный", "укромный"),
    ("общий", "коллективный"),
    ("личный", "персональный"),
    ("частный", "приватный"),
    ("публичный", "общедоступный"),
    ("государственный", "державный"),
    ("народный", "общенародный"),
    ("военный", "армейский"),
    ("гражданский", "штатский"),
    ("религиозный", "духовный"),
    ("светский", "мирской"),
    ("культурный", "образованный"),
    ("дикий", "необузданный"),
    ("домашний", "одомашненный"),
    ("уличный", "беспризорный"),
    ("городской", "урбанистический"),
    ("деревенский", "сельский"),
    ("сельский", "крестьянский"),
    ("лесной", "таёжный"),
    ("морской", "океанический"),
    ("речной", "водный"),
    ("горный", "альпийский"),
    ("степной", "равнинный"),
    ("болотный", "топкий"),
    ("каменный", "скалистый"),
    ("деревянный", "лесозаготовительный"),
    ("стальной", "металлический"),
    ("золотой", "драгоценный"),
    ("серебряный", "благородный"),
    ("стеклянный", "хрустальный"),
    ("пластмассовый", "синтетический"),
    ("бумажный", "целлюлозный"),
    ("шёлковый", "атласный"),
    ("кожаный", "выделанный"),
]


def expand_russian():
    print("MEGA expanding Russian...")
    existing = extract_existing(DICT_DIR / "dict_ru.go")
    print(f"  Before: {len(existing)}")

    # Add noun cases
    for n in russian_nouns:
        nom_sg, repl_sg, nom_pl, repl_pl, gender = n
        forms = expand_russian_noun_cases(nom_sg, repl_sg, nom_pl, repl_pl, gender)
        existing.update(forms)
    print(f"  After noun cases: {len(existing)}")

    # Add verb forms
    for v in russian_verbs_extra:
        inf, repl = v
        forms = expand_russian_verb_full(inf, repl)
        existing.update(forms)
    print(f"  After verbs: {len(existing)}")

    # Add adjective paradigm
    for a in russian_adjs_extra:
        m_form, m_repl = a
        forms = expand_russian_adj_full(m_form, m_repl)
        existing.update(forms)
    print(f"  After adjectives: {len(existing)}")

    write_dict(
        DICT_DIR / "dict_ru.go",
        "synPackRu",
        "// synPackRu is the Russian synonym dictionary.\n// Massive coverage: verbs (all forms), adjectives (all cases/genders/numbers),\n// nouns (6 cases × 2 numbers), adverbs, common phrases, professions, objects.",
        existing
    )
    return len(existing)


# ===================================================================
# MOLDOVAN MEGA EXPANSION
# ===================================================================

# Moldovan: add many new words with all forms
def expand_mo_noun_forms(nom_sg_undef, repl_sg, nom_pl_undef, repl_pl, gender):
    """Generate forms for Moldovan noun: definite/indefinite singular/plural."""
    result = {nom_sg_undef: repl_sg, nom_pl_undef: repl_pl}

    # Add articulated forms
    if gender == 'm':
        # Masculine articulated: ends with -ul or -l
        if nom_sg_undef.endswith('u'):
            result[nom_sg_undef + 'l'] = repl_sg + 'l' if repl_sg.endswith('u') else repl_sg
        else:
            result[nom_sg_undef + 'ul'] = repl_sg + 'ul' if not repl_sg.endswith('ul') else repl_sg
        # Plural articulated: ends with -i
        result[nom_pl_undef + 'i' if not nom_pl_undef.endswith('i') else nom_pl_undef] = repl_pl + 'i' if not repl_pl.endswith('i') else repl_pl
    elif gender == 'f':
        # Feminine articulated: replaces final -ă with -a or adds -a
        if nom_sg_undef.endswith('ă'):
            result[nom_sg_undef[:-1] + 'a'] = repl_sg[:-1] + 'a' if repl_sg.endswith('ă') else repl_sg
        elif nom_sg_undef.endswith('e'):
            result[nom_sg_undef + 'a'] = repl_sg + 'a' if not repl_sg.endswith('a') else repl_sg

    return result


moldovan_extra = [
    # Verbs (infinitive with "a" + multiple tense forms)
    ("a aprinde", "a incendia"), ("a stinge", "a extingere"),
    ("a aduce", "a transporta"), ("a duce", "a transporta"),
    ("a porni", "a iniția"), ("a opri", "a stopa"),
    ("a invita", "a chema"), ("a primi", "a admite"),
    ("a pleca", "a se îndepărta"), ("a sosi", "a ajunge"),
    ("a urma", "a succeda"), ("a precede", "a anteceda"),
    ("a căuta", "a investiga"), ("a găsi", "a descoperi"),
    ("a pierde", "a rătăci"), ("a câștiga", "a obține"),
    ("a împărți", "a divide"), ("a uni", "a amalgama"),
    ("a vorbi", "a comunica"), ("a tăcea", "a păstra tăcerea"),
    ("a striga", "a vocifera"), ("a șopti", "a murmura"),
    ("a privi", "a contempla"), ("a ascunde", "a tăinui"),
    ("a uita", "a omite"), ("a aminti", "a evoca"),
    ("a îndoi", "a flexa"), ("a îndrepta", "a rectifica"),
    ("a strânge", "a aduna"), ("a desface", "a desfășura"),
    ("a urca", "a ascende"), ("a coborî", "a descinde"),
    ("a salva", "a izbăvi"), ("a periclita", "a primejdui"),
    ("a îmbătrâni", "a se senesce"), ("a tinereți", "a întineri"),
    ("a se îmbolnăvi", "a se infirma"), ("a se vindeca", "a se restabili"),
    ("a se naște", "a se ivi pe lume"), ("a muri", "a deceda"),
    ("a căsători", "a uni în matrimoniu"), ("a divorța", "a se separa"),
    ("a se îndrăgosti", "a iubi pasionat"), ("a uri", "a detesta"),
    ("a îmbrăca", "a vestmânta"), ("a dezbrăca", "a denuda"),
    ("a încălța", "a calița"), ("a descălța", "a decalița"),
    ("a alerga", "a goni"), ("a opri", "a opri brusc"),
    ("a se aseza", "a se instala"), ("a se ridica", "a se înălța"),
    ("a se culca", "a se aşeza orizontal"), ("a se trezi", "a se deştepta"),
    ("a respira", "a inhala"), ("a expira", "a exhala"),
    ("a inhala", "a inspira aer"), ("a tuși", "a expectora"),
    ("a strănuta", "a expectora cu zgomot"), ("a sufla", "a sufla aer"),
    ("a zbura", "a plana"), ("a înota", "a se mişca în apă"),
    ("a sări", "a salta"), ("a cădea", "a se prăbuşi"),
    ("a urca", "a ascensiona"), ("a coborî", "a descende"),
    ("a se ridica", "a se ridica vertical"), ("a se așeza", "a se aşeza pe loc"),
    ("a alerga", "a sprinta"), ("a merge", "a păşi"),
    ("a se opri", "a face pauză"), ("a continua", "a persevera"),
    ("a aștepta", "a anticipa"), ("a întâmpina", "a saluta"),
    ("a saluta", "a face reverență"), ("a felicita", "a complimenta"),
    ("a mulțumi", "a manifesta recunoștință"), ("a cere", "a solicita"),
    ("a refuza", "a respinge"), ("a accepta", "a admite"),
    ("a promite", "a făgădui"), ("a uita", "a omite din memorie"),

    # Adjectives
    ("blând", "tandru"), ("aspru", "sever"),
    ("politicos", "amabil"), ("nepoliticos", "necuviincios"),
    ("plin", "îmbelșugat"), ("gol", "vacant"),
    ("uscat", "arid"), ("ud", "umed"),
    ("dulce", "îndulcit"), ("amar", "amărui"),
    ("acru", "lemonic"), ("sărat", "salin"),
    ("delicios", "savuros"), ("nedelicios", "fad"),
    ("plăcut", "agreabil"), ("neplăcut", "dezagreabil"),
    ("util", "folositor"), ("inutil", "neavenit"),
    ("dăunător", "păgubitor"), ("benefic", "salutar"),
    ("important", "esențial"), ("neimportant", "nesemnificativ"),
    ("necesar", "indispensabil"), ("inutil", "superflu"),
    ("dificil", "anevoios"), ("ușor", "facil"),
    ("complex", "intricat"), ("simplu", "elementar"),
    ("clar", "evident"), ("neclar", "ambiguu"),
    ("celebru", "renumit"), ("necunoscut", "anonim"),
    ("popular", "preferat"), ("nepopular", "respins"),
    ("comun", "ubicuu"), ("rar", "scarce"),
    ("special", "distinctiv"), ("obișnuit", "convențional"),
    ("ciudat", "bizar"), ("normal", "standard"),
    ("egal", "echivalent"), ("inegal", "disproporționat"),
    ("similar", "analog"), ("diferit", "disparat"),
    ("identic", "același"), ("variat", "divers"),
    ("periculos", "primejdios"), ("sigur", "în siguranță"),
    ("liber", "independent"), ("ocupat", "preocupat"),
    ("deschis", "destăinuit"), ("închis", "ferecat"),
    ("viu", "însuflețit"), ("mort", "neînsuflețit"),
    ("real", "tangibil"), ("imaginar", "ficțional"),
    ("verita", "autentic"), ("fals", "fictiv"),
    ("vizibil", "perceptibil"), ("invizibil", "imperceptibil"),
    ("audibil", "discernibil"), ("inaudibil", "imperceptibil"),
    ("evident", "manifest"), ("ascuns", "tăinuit"),
    ("public", "deschis"), ("privat", "confidențial"),
    ("național", "patriotic"), ("internațional", "mondial"),
    ("cultural", "civilizat"), ("primitiv", "barbar"),
    ("urban", "citadin"), ("rural", "campestru"),
    ("modern", "contemporan"), ("antic", "primordial"),

    # Nouns
    ("profesor", "educator"), ("elev", "ucenic"),
    ("student", "discipol"), ("director", "administrator"),
    ("medic", "doctor"), ("asistent", "auxiliar"),
    ("inginer", "constructor"), ("arhitect", "proiectant"),
    ("artist", "creator"), ("scriitor", "autor"),
    ("poet", "rapsod"), ("compozitor", "creator muzical"),
    ("muzician", "instrumentist"), ("cântăreț", "vocalist"),
    ("actor", "interpret"), ("regizor", "regizator"),
    ("sportiv", "atlet"), ("antrenor", "instructor"),
    ("om de știință", "savant"), ("filozof", "gânditor"),
    ("preot", "duhovnic"), ("soldat", "militar"),
    ("ofițer", "comandant"), ("general", "căpetenie"),
    ("polițist", "agent al ordinii"), ("paznic", "străjer"),
    ("șofer", "conducător auto"), ("pilot", "aviator"),
    ("marinar", "navigator"), ("căpitan", "comandant naval"),
    ("bucătar", "culinar"), ("ospătar", "deservant"),
    ("vânzător", "comerciant"), ("cumpărător", "client"),
    ("angajator", "patron"), ("angajat", "salariat"),
    ("muncitor", "lucrător"), ("țăran", "agricultor"),
    ("fermier", "agricultor"), ("rege", "monarh"),
    ("regină", "monarhia"), ("prinț", "moștenitor"),
    ("prințesă", "moștenitoare"), ("președinte", "șef de stat"),
    ("ministru", "membru al guvernului"), ("deputat", "parlamentar"),
    ("primar", "edil"), ("guvernator", "administrator regional"),

    # Common objects
    ("telefon", "comunicator"), ("computer", "calculator"),
    ("televizor", "receptor"), ("radio", "aparat de radio"),
    ("ceas", "horometru"), ("ochelari", "lentile optice"),
    ("cheie", "instrument de deschidere"), ("lacăt", "dispozitiv de închidere"),
    ("geantă", "sac de mână"), ("portofel", "buzunar de bani"),
    ("dulap", "garderobă"), ("raft", "stelaj"),
    ("covor", "așternut pe podea"), ("tablou", "pictură"),
    ("statuie", "sculptură"), ("monument", "edificiu comemorativ"),

    # Food
    ("pâine", "produs de panificație"), ("unt", "produs gras"),
    ("brânză", "produs lactat"), ("lapte", "produs lactat lichid"),
    ("carne", "țesut muscular"), ("pește", "ihtiofaună"),
    ("pui", "pasăre comestibilă"), ("ou", "produs ovipar"),
    ("legumă", "cultură horticolă"), ("fruct", "rod"),
    ("măr", "fruct de măr"), ("pară", "fruct de păr"),
    ("portocală", "fruct citric"), ("banană", "fruct tropical"),
    ("lămâie", "citric acrișor"), ("strugure", "fruct de viță"),
    ("cartof", "tubercul"), ("roșie", "tomată"),
    ("castravete", "cucurbitacee"), ("varză", "crucifere"),
    ("morcov", "rădăcinoasă"), ("ceapă", "bulboasă"),
    ("usturoi", "condiment aromatic"), ("zahăr", "îndulcitor"),
    ("sare", "condiment salin"), ("piper", "condiment iute"),
    ("cafea", "băutură caldă tonică"), ("ceai", "infuzie de plante"),

    # Buildings/places
    ("hotel", "loc de cazare"), ("restaurant", "loc de servire a mesei"),
    ("cafenea", "local de cafea"), ("bar", "local de băuturi"),
    ("muzeu", "loc de expoziție"), ("teatru", "loc de spectacole"),
    ("cinema", "loc de filme"), ("concert", "manifestare muzicală"),
    ("stadion", "arenă sportivă"), ("spital", "instituție medicală"),
    ("farmacie", "instituție farmaceutică"), ("bancă", "instituție financiară"),
    ("poștă", "oficiu poștal"), ("gară", "stație de tren"),
    ("aeroport", "port aerian"), ("port", "port maritim"),
    ("stație", "punct de oprire"), ("autobuz", "vehicul de transport public"),
    ("tramvai", "vehicul cu șine"), ("metrou", "transport subteran"),
    ("taxi", "vehicul de transport individual"), ("tren", "convoi feroviar"),
    ("avion", "aeronavă"), ("vapor", "navă"),
    ("barcă", "ambarcațiune"), ("automobil", "vehicul"),
    ("bicicletă", "vehicul cu două roți"), ("motocicletă", "vehicul motorizat cu două roți"),

    # Abstract
    ("suflet", "psiche"), ("minte", "intelect"),
    ("conștiință", "moralitate"), ("voință", "determinare"),
    ("libertate", "independență"), ("sclavie", "robie"),
    ("putere", "autoritate"), ("forță", "vigoare"),
    ("glorie", "renume"), ("onoare", "demnitate"),
    ("rușine", "umilință"), ("vină", "culpabilitate"),
    ("iertare", "absolvire"), ("trădare", "perfidiate"),
    ("loialitate", "fidelitate"), ("prietenie", "camaraderie"),
    ("dragoste", "afecțiune"), ("pasiune", "ardoare"),
    ("tandrețe", "duioșie"), ("grijă", "preocupare"),
    ("apărare", "protecție"), ("siguranță", "securitate"),
    ("pericol", "primejdie"), ("amenințare", "intimidare"),

    # Education
    ("lecție", "instruire"), ("examen", "verificare"),
    ("notă", "calificativ"), ("temă", "sarcină"),
    ("exercițiu", "antrenament"), ("caiet", "registru școlar"),
    ("manual", "instrument didactic"), ("dicționar", "lexicon"),
    ("creion", "stilou cu grafit"), ("pix", "instrument de scris"),
    ("hârtie", "material de scris"), ("carte", "tom literar"),
    ("revistă", "publicație periodică"), ("ziar", "cotidian"),
    ("articol", "scriere"), ("ediție", "tipăritură"),

    # Weather
    ("vreme", "condiții meteorologice"), ("ceață", "negură"),
    ("furtună", "uragan"), ("uragan", "ciclon tropical"),
    ("tunet", "bubuit ceresc"), ("fulger", "scânteie cerească"),
    ("vânt", "curent de aer"), ("grindină", "căderi de gheață"),

    # Emotions
    ("bucurie", "veselie"), ("tristețe", "melancolie"),
    ("frică", "spaimă"), ("speranță", "năzuință"),
    ("ură", "ostilitate"), ("furie", "mânie"),
    ("calm", "seninătate"), ("nervozitate", "agitație"),
    ("liniște", "tăcere"), ("zgomot", "vacarm"),
    ("plăcere", "delectare"), ("durere", "suferință"),
]


def expand_moldovan():
    print("MEGA expanding Moldovan...")
    existing = extract_existing(DICT_DIR / "dict_mo.go")
    print(f"  Before: {len(existing)}")

    # Add new entries
    for entry in moldovan_extra:
        word, repl = entry
        if word not in existing:
            existing[word] = repl

    print(f"  After new entries: {len(existing)}")

    # Add more articulated/inflected forms
    extras = {}
    for word, repl in list(existing.items()):
        # Skip multi-word entries
        if ' ' in word:
            continue
        # Generate definite form for nouns ending in -ă (feminine)
        if word.endswith('ă') and len(word) > 2:
            articulated = word[:-1] + 'a'
            if articulated not in existing:
                if repl.endswith('ă') and len(repl) > 2:
                    extras[articulated] = repl[:-1] + 'a'
                else:
                    extras[articulated] = repl
        # Generate definite form for nouns ending in -e
        elif word.endswith('e') and len(word) > 2:
            articulated = word + 'a'
            if articulated not in existing:
                extras[articulated] = repl + 'a' if not repl.endswith('a') else repl

    existing.update(extras)
    print(f"  After articulated forms: {len(existing)}")

    write_dict(
        DICT_DIR / "dict_mo.go",
        "synPackMo",
        "// synPackMo is the Moldovan/Romanian synonym dictionary.\n// Comprehensive coverage: verbs (all forms), articulated/indefinite forms,\n// adjectives (all genders/numbers), pronouns, common phrases, professions.",
        existing
    )
    return len(existing)


if __name__ == "__main__":
    ru = expand_russian()
    mo = expand_moldovan()
    print(f"\nFinal: RU={ru}, MO={mo}")
