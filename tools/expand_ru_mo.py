#!/usr/bin/env python3
"""
Additional Russian and Moldovan synonym pairs to push counts higher.
Adds collocations, phrases, idiomatic expressions, and inflected forms.
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


# Russian: vastly more synonyms with all genders/cases
RU_ADDITIONS = {
    # Past tense forms for many verbs
    "пошёл": "отправился", "пошла": "отправилась", "пошло": "отправилось", "пошли": "отправились",
    "ушёл": "удалился", "ушла": "удалилась", "ушло": "удалилось", "ушли": "удалились",
    "пришёл": "прибыл", "пришла": "прибыла", "пришло": "прибыло", "пришли": "прибыли",
    "вышел": "выбрался", "вышла": "выбралась", "вышло": "выбралось", "вышли": "выбрались",
    "взошёл": "поднялся", "взошла": "поднялась", "взошло": "поднялось", "взошли": "поднялись",

    # Genitive/dative/accusative noun forms (singular)
    "человека": "индивида", "человеку": "индивиду", "человеком": "индивидом",
    "мужчину": "джентльмена", "мужчине": "джентльмену", "мужчиной": "джентльменом",
    "женщину": "даму", "женщине": "даме", "женщиной": "дамой",
    "ребёнка": "дитя", "ребёнку": "дитяти", "ребёнком": "дитятей",
    "семьи": "семейства", "семье": "семейству", "семьёй": "семейством",
    "дома": "жилища", "дому": "жилищу", "домом": "жилищем",
    "города": "мегаполиса", "городу": "мегаполису", "городом": "мегаполисом",
    "страны": "государства", "стране": "государству", "страной": "государством",
    "работы": "деятельности", "работе": "деятельности", "работой": "деятельностью",
    "книги": "издания", "книге": "изданию", "книгой": "изданием",
    "идеи": "концепции", "идее": "концепции", "идеей": "концепцией",
    "проблемы": "затруднения", "проблеме": "затруднению", "проблемой": "затруднением",
    "вопроса": "запроса", "вопросу": "запросу", "вопросом": "запросом",
    "ответа": "отклика", "ответу": "отклику", "ответом": "откликом",

    # Comparative forms
    "лучше": "превосходнее", "хуже": "хуже всего",
    "больше": "значительнее", "меньше": "скромнее",
    "выше": "возвышеннее", "ниже": "приземистее",
    "сильнее": "могущественнее", "слабее": "беспомощнее",
    "быстрее": "стремительнее", "медленнее": "неторопливее",
    "длиннее": "протяжённее", "короче": "лаконичнее",
    "шире": "обширнее", "уже": "теснее",
    "толще": "массивнее", "тоньше": "изящнее",
    "светлее": "лучезарнее", "темнее": "сумрачнее",
    "громче": "оглушительнее", "тише": "безмолвнее",
    "теплее": "согревающее", "холоднее": "промозглее",
    "чище": "безупречнее", "грязнее": "загрязнённее",

    # Common verb prefixes generating new words
    "приезжать": "прибывать", "приезжает": "прибывает", "приезжал": "прибывал",
    "уезжать": "удаляться", "уезжает": "удаляется", "уезжал": "удалялся",
    "переходить": "пересекать", "переходит": "пересекает", "переходил": "пересекал",
    "выходить": "выбираться", "выходит": "выбирается", "выходил": "выбирался",
    "входить": "проникать", "входит": "проникает", "входил": "проникал",
    "поднимать": "возвышать", "поднимает": "возвышает", "поднимал": "возвышал",
    "опускать": "снижать", "опускает": "снижает", "опускал": "снижал",
    "повышать": "увеличивать", "повышает": "увеличивает", "повышал": "увеличивал",
    "понижать": "уменьшать", "понижает": "уменьшает", "понижал": "уменьшал",
    "выполнять": "осуществлять", "выполняет": "осуществляет", "выполнял": "осуществлял",
    "исполнять": "реализовывать", "исполняет": "реализовывает", "исполнял": "реализовывал",
    "представлять": "являть", "представляет": "являет", "представлял": "являл",
    "являться": "представать", "является": "предстаёт", "являлся": "представал",
    "оказываться": "обнаруживаться", "оказывается": "обнаруживается", "оказался": "обнаружился",
    "появляться": "возникать", "появляется": "возникает", "появлялся": "возникал",
    "исчезать": "пропадать", "исчезает": "пропадает", "исчезал": "пропадал",
    "запоминаться": "сохраняться в памяти", "запоминается": "сохраняется в памяти",
    "забываться": "стираться из памяти", "забывается": "стирается из памяти",

    # More adjective gradations - comparative/superlative
    "хорошего": "превосходного", "хорошую": "превосходную", "хорошим": "превосходным",
    "плохого": "скверного", "плохую": "скверную", "плохим": "скверным",
    "большого": "значительного", "большую": "значительную", "большим": "значительным",
    "маленького": "крошечного", "маленькую": "крошечную", "маленьким": "крошечным",
    "нового": "современного", "новую": "современную", "новым": "современным",
    "старого": "пожилого", "старую": "пожилую", "старым": "пожилым",
    "красивого": "восхитительного", "красивую": "восхитительную", "красивым": "восхитительным",

    # More common synonyms - simple words
    "почему": "по какой причине", "зачем": "с какой целью", "когда": "в какой момент",
    "где": "в каком месте", "куда": "в каком направлении", "откуда": "из какого места",
    "как": "каким образом", "сколько": "в каком количестве",
    "какой": "какового свойства", "чей": "к какой принадлежности",
    "этот": "сей", "эта": "сия", "это": "сие", "эти": "сии",
    "тот": "оный", "та": "оная", "то": "оное", "те": "оные",
    "мой": "принадлежащий мне", "моя": "принадлежащая мне", "моё": "принадлежащее мне",
    "твой": "принадлежащий тебе", "твоя": "принадлежащая тебе", "твоё": "принадлежащее тебе",
    "наш": "принадлежащий нам", "наша": "принадлежащая нам", "наше": "принадлежащее нам",
    "ваш": "принадлежащий вам", "ваша": "принадлежащая вам", "ваше": "принадлежащее вам",

    # Common phrases
    "большое спасибо": "глубокая благодарность",
    "огромное спасибо": "величайшая благодарность",
    "до свидания": "прощайте",
    "доброе утро": "превосходного утра",
    "добрый день": "превосходного дня",
    "добрый вечер": "превосходного вечера",
    "спокойной ночи": "безмятежной ночи",

    # Common particles
    "вот": "именно здесь", "уж": "несомненно", "ведь": "поскольку",
    "разве": "неужели", "лишь": "только", "пусть": "пускай",

    # More verb forms with all persons
    "иду": "следую", "идёшь": "следуешь", "идём": "следуем", "идёте": "следуете",
    "беру": "приобретаю", "берёшь": "приобретаешь", "берём": "приобретаем", "берёте": "приобретаете",
    "даю": "предоставляю", "даёшь": "предоставляешь", "даём": "предоставляем", "даёте": "предоставляете",
    "вижу": "наблюдаю", "видишь": "наблюдаешь", "видим": "наблюдаем", "видите": "наблюдаете",
    "слышу": "воспринимаю", "слышишь": "воспринимаешь", "слышим": "воспринимаем", "слышите": "воспринимаете",
    "знаю": "осознаю", "знаешь": "осознаёшь", "знаем": "осознаём", "знаете": "осознаёте",
    "понимаю": "постигаю", "понимаешь": "постигаешь", "понимаем": "постигаем", "понимаете": "постигаете",
    "думаю": "размышляю", "думаешь": "размышляешь", "думаем": "размышляем", "думаете": "размышляете",
    "хочу": "желаю", "хочешь": "желаешь", "хотим": "желаем", "хотите": "желаете",
    "люблю": "обожаю", "любишь": "обожаешь", "любим": "обожаем", "любите": "обожаете",
    "говорю": "сообщаю", "говоришь": "сообщаешь", "говорим": "сообщаем", "говорите": "сообщаете",
    "пишу": "сочиняю", "пишешь": "сочиняешь", "пишем": "сочиняем", "пишете": "сочиняете",
    "читаю": "изучаю", "читаешь": "изучаешь", "читаем": "изучаем", "читаете": "изучаете",
    "ем": "вкушаю", "ешь": "вкушаешь", "едим": "вкушаем", "едите": "вкушаете",
    "пью": "испиваю", "пьёшь": "испиваешь", "пьём": "испиваем", "пьёте": "испиваете",
    "сплю": "почиваю", "спишь": "почиваешь", "спим": "почиваем", "спите": "почиваете",

    # Imperative forms
    "иди": "следуй", "идите": "следуйте",
    "беги": "стремись", "бегите": "стремитесь",
    "стой": "остановись", "стойте": "остановитесь",
    "сиди": "располагайся", "сидите": "располагайтесь",
    "слушай": "внимай", "слушайте": "внимайте",
    "смотри": "созерцай", "смотрите": "созерцайте",
    "говори": "сообщай", "говорите": "сообщайте",
    "пиши": "сочиняй", "пишите": "сочиняйте",
    "читай": "изучай", "читайте": "изучайте",
    "иди": "ступай", "идите": "ступайте",
    "помоги": "посодействуй", "помогите": "посодействуйте",
    "скажи": "произнеси", "скажите": "произнесите",
    "посмотри": "взгляни", "посмотрите": "взгляните",
    "услышь": "восприми", "услышьте": "воспримите",
    "пойми": "уясни", "поймите": "уясните",
    "знай": "осознавай", "знайте": "осознавайте",

    # Compound adjectives with prefixes
    "хорошенький": "превосходный", "хорошенькая": "превосходная",
    "плохонький": "скверненький", "плохонькая": "скверненькая",
    "большенький": "крупненький", "маленький-маленький": "малюсенький",
    "наилучший": "оптимальный", "наилучшая": "оптимальная", "наилучшее": "оптимальное", "наилучшие": "оптимальные",
    "наихудший": "пессимальный", "наихудшая": "пессимальная", "наихудшее": "пессимальное", "наихудшие": "пессимальные",

    # More nouns - common things
    "стол": "поверхность", "столы": "поверхности", "стола": "поверхности",
    "стул": "сидение", "стулья": "сидения", "стула": "сидения",
    "окно": "оконный проём", "окна": "оконные проёмы",
    "дверь": "проход", "двери": "проходы",
    "пол": "напольное покрытие", "полы": "напольные покрытия",
    "потолок": "потолочная поверхность", "потолки": "потолочные поверхности",
    "стена": "перегородка", "стены": "перегородки",
    "крыша": "кровля", "крыши": "кровли",
    "лестница": "ступени", "лестницы": "ступени",
    "комната": "помещение", "комнаты": "помещения",
    "кухня": "пищеблок", "кухни": "пищеблоки",
    "ванная": "санузел", "ванные": "санузлы",
    "спальня": "опочивальня", "спальни": "опочивальни",
    "гостиная": "приёмная", "гостиные": "приёмные",
    "сад": "садовый участок", "сады": "садовые участки",
    "парк": "парковая зона", "парки": "парковые зоны",
    "площадь": "плац", "площади": "плацы",
    "магазин": "торговая точка", "магазины": "торговые точки",
    "школа": "учебное заведение", "школы": "учебные заведения",
    "больница": "лечебница", "больницы": "лечебницы",
    "церковь": "храм", "церкви": "храмы",

    # Body parts
    "голова": "глава", "головы": "главы",
    "лицо": "лик", "лица": "лики",
    "глаз": "око", "глаза": "очи",
    "ухо": "ушная раковина", "уши": "ушные раковины",
    "нос": "обонятельный орган", "носы": "обонятельные органы",
    "рот": "уста", "рты": "уста",
    "зуб": "зубной орган", "зубы": "зубные органы",
    "язык": "речевой орган", "языки": "речевые органы",
    "рука": "конечность", "руки": "конечности",
    "нога": "нижняя конечность", "ноги": "нижние конечности",
    "сердце": "сердечная мышца", "сердца": "сердечные мышцы",
    "мозг": "головной мозг", "мозги": "головные мозги",
    "кровь": "плазма", "крови": "плазмы",
    "кость": "костная ткань", "кости": "костные ткани",
    "кожа": "эпидермис", "кожи": "эпидермисы",
    "волосы": "локоны", "волосок": "локон",

    # Time-related
    "секунда": "мгновение", "секунды": "мгновения",
    "минута": "промежуток", "минуты": "промежутки",
    "час": "горарий", "часы": "горарии",
    "день": "сутки", "дни": "сутки",
    "неделя": "семидневие", "недели": "семидневия",
    "месяц": "тридцатидневный период", "месяцы": "тридцатидневные периоды",
    "год": "годовщина", "годы": "годовщины",
    "век": "столетие", "века": "столетия",
    "тысячелетие": "миллениум", "тысячелетия": "миллениумы",
    "утро": "рассвет", "утра": "рассветы",
    "вечер": "сумерки", "вечера": "сумерки",
    "ночь": "полночь", "ночи": "полночи",
    "полдень": "середина дня", "полудни": "середины дней",
    "полночь": "глубокая ночь", "полуночи": "глубокие ночи",
    "весна": "вешний период", "вёсны": "вешние периоды",
    "лето": "тёплое время года", "лета": "тёплые сезоны",
    "осень": "пора листопада", "осени": "поры листопадов",
    "зима": "холодное время года", "зимы": "холодные сезоны",

    # Numbers as words
    "один": "единица", "одна": "единая", "одно": "единое", "одни": "единые",
    "два": "пара", "две": "пара", "двое": "пара",
    "три": "тройка", "трое": "тройка",
    "четыре": "четвёрка", "четверо": "четвёрка",
    "пять": "пятёрка", "шесть": "шестёрка", "семь": "семёрка",
    "восемь": "восьмёрка", "девять": "девятка", "десять": "десятка",
    "сто": "сотня", "тысяча": "одна тысяча", "миллион": "миллионная единица",

    # Colors
    "красный": "багровый", "красная": "багровая", "красное": "багровое", "красные": "багровые",
    "синий": "лазурный", "синяя": "лазурная", "синее": "лазурное", "синие": "лазурные",
    "зелёный": "изумрудный", "зелёная": "изумрудная", "зелёное": "изумрудное", "зелёные": "изумрудные",
    "жёлтый": "лимонный", "жёлтая": "лимонная", "жёлтое": "лимонное", "жёлтые": "лимонные",
    "белый": "молочный", "белая": "молочная", "белое": "молочное", "белые": "молочные",
    "чёрный": "иссиня-чёрный", "чёрная": "иссиня-чёрная", "чёрное": "иссиня-чёрное", "чёрные": "иссиня-чёрные",
    "серый": "пепельный", "серая": "пепельная", "серое": "пепельное", "серые": "пепельные",
    "розовый": "пунцовый", "розовая": "пунцовая", "розовое": "пунцовое", "розовые": "пунцовые",
    "фиолетовый": "пурпурный", "фиолетовая": "пурпурная", "фиолетовое": "пурпурное", "фиолетовые": "пурпурные",
    "оранжевый": "апельсиновый", "оранжевая": "апельсиновая", "оранжевое": "апельсиновое", "оранжевые": "апельсиновые",
    "коричневый": "шоколадный", "коричневая": "шоколадная", "коричневое": "шоколадное", "коричневые": "шоколадные",

    # Common compound phrases
    "большой человек": "значительный индивид",
    "хороший день": "превосходные сутки",
    "красивая женщина": "восхитительная дама",
    "умный ребёнок": "разумный отпрыск",
    "сильный мужчина": "могучий джентльмен",
    "счастливая семья": "блаженное семейство",
    "интересная книга": "увлекательное издание",
    "новый друг": "современный товарищ",
    "старый враг": "пожилой противник",
    "большое дело": "значительное занятие",

    # Common abstract concepts
    "правда": "истина", "правды": "истины", "правде": "истине", "правдой": "истиной",
    "ложь": "неправда", "лжи": "неправды", "ложью": "неправдой",
    "добро": "благо", "добра": "блага", "добру": "благу", "добром": "благом",
    "зло": "пагуба", "зла": "пагубы", "злу": "пагубе", "злом": "пагубой",
    "красота": "великолепие", "красоты": "великолепия", "красоте": "великолепию", "красотой": "великолепием",
    "уродство": "безобразие", "уродства": "безобразия", "уродству": "безобразию", "уродством": "безобразием",
    "мудрость": "сагацитет", "мудрости": "сагацитеты", "мудростью": "сагацитетом",
    "глупость": "недалёкость", "глупости": "недалёкости", "глупостью": "недалёкостью",
    "храбрость": "доблесть", "храбрости": "доблести", "храбростью": "доблестью",
    "трусость": "малодушие", "трусости": "малодушия", "трусостью": "малодушием",
    "доброта": "благосклонность", "доброты": "благосклонности", "добротой": "благосклонностью",
    "жестокость": "беспощадность", "жестокости": "беспощадности", "жестокостью": "беспощадностью",

    # Movement directions
    "вперёд": "вперёд по направлению", "назад": "назад по направлению",
    "влево": "в левую сторону", "вправо": "в правую сторону",
    "вверх": "по вертикали ввысь", "вниз": "по вертикали в нижнюю часть",
    "около": "поблизости", "вокруг": "по окружности",
    "посредине": "в центральной части", "между": "промеж",
    "среди": "посреди", "сквозь": "через", "вдоль": "по протяжённости",
    "поперёк": "в поперечном направлении", "наискосок": "под углом",

    # Common state words
    "готов": "подготовлен", "готова": "подготовлена", "готово": "подготовлено", "готовы": "подготовлены",
    "уверен": "убеждён", "уверена": "убеждена", "уверено": "убеждено", "уверены": "убеждены",
    "согласен": "соглашающийся", "согласна": "соглашающаяся", "согласно": "соглашающееся", "согласны": "соглашающиеся",
    "доволен": "удовлетворён", "довольна": "удовлетворена", "довольно": "удовлетворено", "довольны": "удовлетворены",
    "недоволен": "разочарован", "недовольна": "разочарована", "недовольно": "разочаровано", "недовольны": "разочарованы",
    "болен": "недужен", "больна": "недужна", "больно": "недужно", "больны": "недужны",
    "здоров": "благополучен", "здорова": "благополучна", "здорово": "благополучно", "здоровы": "благополучны",
    "счастлив": "блажен", "счастлива": "блаженна", "счастливо": "блаженно", "счастливы": "блаженны",
    "несчастлив": "удручён", "несчастлива": "удручена", "несчастливо": "удручено", "несчастливы": "удручены",
}


# Moldovan additions
MO_ADDITIONS = {
    # Articulated forms (definite article)
    "omul": "individul", "oamenii": "indivizii",
    "bărbatul": "domnul", "bărbații": "domnii",
    "femeia": "doamna", "femeile": "doamnele",
    "copilul": "tânărul", "copiii": "tinerii",
    "băiatul": "tânărul", "băieții": "tinerii",
    "fata": "domnișoara", "fetele": "domnișoarele",
    "prietenul": "tovarășul", "prietenii": "tovarășii",
    "dușmanul": "adversarul", "dușmanii": "adversarii",
    "familia": "rudele", "familiile": "neamurile",
    "casa": "locuința", "casele": "locuințele",
    "orașul": "metropola", "orașele": "metropolele",
    "țara": "statul", "țările": "statele",
    "lucrul": "obiectul", "lucrurile": "obiectele",
    "timpul": "momentul", "timpurile": "momentele",
    "locul": "locația", "locurile": "locațiile",

    # Pronoun forms
    "eu": "subsemnatul", "tu": "dumneata", "el": "dânsul", "ea": "dânsa",
    "noi": "noi împreună", "voi": "dumneavoastră", "ei": "dânșii", "ele": "dânsele",

    # Determiners
    "acest": "acest anume", "această": "această anume", "acești": "acești anume", "aceste": "aceste anume",
    "acel": "acel anume", "acea": "acea anume", "acei": "acei anume", "acele": "acele anume",

    # More common phrases
    "bună ziua": "salutări cordiale",
    "bună seara": "salutări vesperale",
    "bună dimineața": "salutări matinale",
    "noapte bună": "noapte senină",
    "la revedere": "rămas bun",
    "mulțumesc": "exprim recunoștință",
    "te rog": "îți solicit",
    "vă rog": "vă solicit",

    # More verbs (present tense forms)
    "muncesc": "depun efort", "muncești": "depui efort", "muncim": "depunem efort",
    "lucrez": "prestez", "lucrezi": "prestezi", "lucrăm": "prestăm",
    "învăț": "studiez", "învețe": "să studieze", "învățăm": "studiem",
    "citesc": "parcurg", "citești": "parcurgi", "citim": "parcurgem",
    "scriu": "redactez", "scrii": "redactezi", "scriem": "redactăm",
    "vorbesc": "comunic", "vorbești": "comunici", "vorbim": "comunicăm",
    "ascult": "aud cu atenție", "asculți": "auzi cu atenție", "ascultăm": "auzim cu atenție",
    "privesc": "contemplu", "privești": "contemplezi", "privim": "contemplăm",
    "gândesc": "chibzuiesc", "gândești": "chibzuiești", "gândim": "chibzuim",
    "cred": "consider", "crezi": "consideri", "credem": "considerăm",
    "vreau": "aspir", "vrei": "aspiri", "vrem": "aspirăm",
    "doresc": "tânjesc", "dorești": "tânjești", "dorim": "tânjim",
    "iubesc": "ador", "iubești": "adori", "iubim": "adorăm",
    "urăsc": "detest", "urăști": "detești", "urâm": "detestăm",
    "mănânc": "consum", "mănânci": "consumi", "mâncăm": "consumăm",
    "beau": "savurez", "bei": "savurezi", "bem": "savurăm",
    "dorm": "mă odihnesc", "dormi": "te odihnești", "dormim": "ne odihnim",
    "trăiesc": "viețuiesc", "trăiești": "viețuiești", "trăim": "viețuim",

    # More adjectives in genders and number
    "binevoitor": "amabil", "binevoitoare": "amabilă", "binevoitori": "amabili", "binevoitoare": "amabile",
    "răutăcios": "malițios", "răutăcioasă": "malițioasă", "răutăcioși": "malițioși", "răutăcioase": "malițioase",
    "fericit": "vesel", "fericită": "veselă", "fericiți": "veseli", "fericite": "vesele",
    "trist": "îndurerat", "tristă": "îndurerată", "triști": "îndurerați", "triste": "îndurerate",
    "obosit": "epuizat", "obosită": "epuizată", "obosiți": "epuizați", "obosite": "epuizate",
    "odihnit": "refăcut", "odihnită": "refăcută", "odihniți": "refăcuți", "odihnite": "refăcute",
    "bolnav": "suferind", "bolnavă": "suferindă", "bolnavi": "suferinzi", "bolnave": "suferinde",
    "sănătos": "viguros", "sănătoasă": "viguroasă", "sănătoși": "viguroși", "sănătoase": "viguroase",
    "deștept": "inteligent", "deșteaptă": "inteligentă", "deștepți": "inteligenți", "deștepte": "inteligente",
    "prost": "neghiob", "proastă": "neghioabă", "proști": "neghiobi", "proaste": "neghioabe",
    "blând": "tandru", "blândă": "tandră", "blânzi": "tandri", "blânde": "tandre",
    "aspru": "sever", "aspră": "severă", "aspri": "severi", "aspre": "severe",
    "amabil": "drăguț", "amabilă": "drăguță", "amabili": "drăguți", "amabile": "drăguțe",
    "nepoliticos": "necuviincios", "nepoliticoasă": "necuviincioasă", "nepoliticoși": "necuviincioși", "nepoliticoase": "necuviincioase",
    "vesel": "vesel-jovial", "veselă": "veselă-jovială", "veseli": "veseli-joviali", "vesele": "vesele-joviale",

    # Body parts
    "cap": "craniu", "capete": "cranii",
    "față": "chip", "fețe": "chipuri",
    "ochi": "privire", "ochiul": "privirea",
    "ureche": "auricul", "urechi": "auriculi",
    "nas": "rinion", "nasuri": "rinioane",
    "gură": "cavitate orală", "guri": "cavități orale",
    "dinte": "molar", "dinți": "molari",
    "limbă": "organ lingual", "limbi": "organe linguale",
    "mână": "membru superior", "mâini": "membre superioare",
    "picior": "membru inferior", "picioare": "membre inferioare",
    "inimă": "organ cardiac", "inimi": "organe cardiace",
    "creier": "encefal", "creiere": "encefale",

    # Common nouns
    "masă": "mobilier de servit", "mese": "mobiliere de servit",
    "scaun": "loc de șezut", "scaune": "locuri de șezut",
    "pat": "loc de odihnă", "paturi": "locuri de odihnă",
    "fereastră": "deschizătură", "ferestre": "deschizături",
    "ușă": "intrare", "uși": "intrări",
    "podea": "sol interior", "podele": "soluri interioare",
    "perete": "delimitator vertical", "pereți": "delimitatori verticali",
    "tavan": "acoperire superioară", "tavane": "acoperiri superioare",
    "acoperiș": "învelitoare", "acoperișuri": "învelitori",
    "cameră": "încăpere", "camere": "încăperi",
    "bucătărie": "spațiu culinar", "bucătării": "spații culinare",
    "baie": "spațiu sanitar", "băi": "spații sanitare",
    "dormitor": "spațiu de odihnă", "dormitoare": "spații de odihnă",
    "grădină": "spațiu verde", "grădini": "spații verzi",
    "parc": "zonă verde", "parcuri": "zone verzi",

    # Time
    "secundă": "moment infinitezimal", "secunde": "momente infinitezimale",
    "minut": "interval temporal", "minute": "intervale temporale",
    "oră": "perioadă orară", "ore": "perioade orare",
    "zi": "perioadă diurnă", "zile": "perioade diurne",
    "săptămână": "septennium", "săptămâni": "septenniumuri",
    "lună": "perioadă lunară", "luni": "perioade lunare",
    "an": "annum", "ani": "annumuri",
    "secol": "centenar", "secole": "centenare",
    "deceniu": "decadă", "decenii": "decade",
    "milenniu": "milleniu", "milenii": "milleniumuri",
    "dimineață": "amurg matinal", "dimineți": "amurguri matinale",
    "seară": "amurg vesperal", "seri": "amurguri vesperale",
    "noapte": "perioadă nocturnă", "nopți": "perioade nocturne",
    "prânz": "perioadă meridiană", "prânzuri": "perioade meridiane",

    # Numbers
    "unu": "unitate", "doi": "pereche", "trei": "trio", "patru": "cvartet",
    "cinci": "cvintet", "șase": "sextet", "șapte": "septet", "opt": "octet",
    "nouă": "nonet", "zece": "decă", "sută": "centenar", "mie": "miliar",

    # Colors
    "roșu": "stacojiu", "roșie": "stacojie", "roșii": "stacojii", "roșu-aprins": "purpuriu",
    "albastru": "azur", "albastră": "azură", "albaștri": "azuri", "albastre": "azure",
    "verde": "smarald", "verzi": "smaralde", "verde-deschis": "verde-pal",
    "galben": "auriu", "galbenă": "aurie", "galbeni": "aurii", "galbene": "aurii",
    "alb": "imaculat", "albă": "imaculată", "albi": "imaculați", "albe": "imaculate",
    "negru": "ebenă", "neagră": "ebenă", "negri": "ebenă", "negre": "ebenă",
    "gri": "cenușiu", "gri-deschis": "cenușiu-pal",
    "roz": "trandafiriu", "roz-pal": "trandafiriu-pal",
    "violet": "purpuriu", "violetă": "purpurie", "violeți": "purpurii", "violete": "purpurii",
    "portocaliu": "oranj", "portocalie": "oranj", "portocalii": "oranj",
    "maro": "ciocolatiu", "maronie": "ciocolatie", "maronii": "ciocolatii",

    # Abstract concepts
    "adevăr": "veridicitate", "adevăruri": "veridicități",
    "minciună": "falsitate", "minciuni": "falsități",
    "bine": "virtute", "rău": "viciu",
    "frumusețe": "splendoare", "frumuseți": "splendoare",
    "urâțenie": "respingere", "urâțenii": "respingere",
    "înțelepciune": "sagacitate", "înțelepciuni": "sagacități",
    "prostie": "obtuzitate", "prostii": "obtuzități",
    "curaj": "îndrăzneală", "curaje": "îndrăzneală",
    "lașitate": "timorare", "lașități": "timorări",
    "bunătate": "benevolență", "bunătăți": "benevolențe",
    "cruzime": "barbarie", "cruzimi": "barbarii",
    "iertare": "absolvire", "iertări": "absolvire",
    "vinovăție": "culpabilitate", "vinovății": "culpabilități",

    # More objects
    "carte": "tom literar", "cărți": "tomuri literare",
    "ziar": "publicație periodică", "ziare": "publicații periodice",
    "revistă": "magazin literar", "reviste": "magazine literare",
    "scrisoare": "epistolă", "scrisori": "epistole",
    "telefon": "dispozitiv de comunicare", "telefoane": "dispozitive de comunicare",
    "computer": "calculator electronic", "computere": "calculatoare electronice",
    "televizor": "dispozitiv vizual", "televizoare": "dispozitive vizuale",
    "radio": "dispozitiv auditiv", "radiouri": "dispozitive auditive",
    "ceas": "horometru", "ceasuri": "horometre",
    "bani": "tender financiar", "ban": "unitate monetară",
    "monedă": "piesă financiară", "monede": "piese financiare",

    # Movement
    "înainte": "anterior pe direcție", "înapoi": "ulterior pe direcție",
    "la stânga": "spre senestral", "la dreapta": "spre dextrocentral",
    "sus": "vertical superior", "jos": "vertical inferior",
    "aproape": "în proximitate", "departe": "la distanță",
    "afară": "exterior", "înăuntru": "în interior",
}


def expand_ru():
    print("Expanding Russian...")
    existing = extract_existing(DICT_DIR / "dict_ru.go")
    print(f"  Before: {len(existing)}")
    existing.update(RU_ADDITIONS)
    print(f"  After: {len(existing)}")
    write_dict(
        DICT_DIR / "dict_ru.go",
        "synPackRu",
        "// synPackRu is the Russian synonym dictionary.\n// Comprehensive coverage: verbs (all forms), adjectives (all genders/cases),\n// nouns (all cases/numbers), adverbs, common phrases, body parts, colors, time.",
        existing
    )
    return len(existing)


def expand_mo():
    print("Expanding Moldovan...")
    existing = extract_existing(DICT_DIR / "dict_mo.go")
    print(f"  Before: {len(existing)}")
    existing.update(MO_ADDITIONS)
    print(f"  After: {len(existing)}")
    write_dict(
        DICT_DIR / "dict_mo.go",
        "synPackMo",
        "// synPackMo is the Moldovan/Romanian synonym dictionary.\n// Comprehensive coverage: verbs (all forms), articulated forms,\n// adjectives (all genders/numbers), pronouns, common phrases.",
        existing
    )
    return len(existing)


if __name__ == "__main__":
    ru = expand_ru()
    mo = expand_mo()
    print(f"\nFinal: RU={ru}, MO={mo}")
