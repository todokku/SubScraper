let netflixInit = false;
let primeInit = false;
let hostUrl = location.href;
if (hostUrl.includes("netflix") || hostUrl.includes("youtube")) {
    netflixInit = false;
    startListener(hostUrl);
} else if (hostUrl.includes("primevideo")) {
    primeInit = false;
    startListener(hostUrl);
}

function getMeaning(source, video) {
    let val = [];
    var subtitles;
    if (source === "netflix") {
        subtitles = document.querySelector(".player-timedtext");
    } else if (source === "primevideo") {
        subtitles = document.querySelector('.persistentPanel');
    } else if (source === "youtube") {
        subtitles = document.querySelector(".captions-text");
    }

    if (subtitles && video.paused) {
        if (subtitles.firstChild != null && val.length === 0) {
            let children = subtitles.children;
            for (let i = 0; i < children.length; i++) {
                filteredWords = children[i].innerText.removeStopWords();
                let x = [];
                if (source === "primevideo") {
                    x = filteredWords.replace(/(\r\n|\n|\r)/gm," ").split(" ")
                } else {
                    x = filteredWords.split(" ");
                }
                x.forEach((e) => {
                    let y = e.replace(/[^\w\s]/gi, '');
                    val.push(y)
                });
                //removing duplicatesu
                val = val.filter(function(item,i,allItems){
                    return i===allItems.indexOf(item);
                });
            }
            let count = 0;
            let listHtml = '<ul class="list-group">';
            val.forEach((v) => {
                let request = new XMLHttpRequest();
                request.open('GET', "http://localhost:8000/api/meaning/" + v, true);
                request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
                request.setRequestHeader('Accept', 'application/json');
                request.send();

                request.onreadystatechange = function () {
                    if (request.readyState === 4 && request.status === 200) {
                        let resp = JSON.parse(request.responseText);
                        if (resp.includes("Error 404")) {
                            console.log("Invalid Word")
                        } else {
                            let obj = JSON.parse(resp);
                            if (obj.hasOwnProperty('definitions')) {
                                listHtml += '<li class="list-group-item">' + v + ": " + obj.definitions[0].definition + '</li>';
                            } else {
                                listHtml += '<li class="list-group-item">' + v + ': No meaning</li>';
                            }
                        }
                        count++
                    }
                };
            });
            listHtml += '</ul>';
            let intervalId = setInterval(() => {
                if (count === val.length) {
                    Swal.fire({
                        position: 'top-end',
                        title: 'Word List',
                        html: listHtml,
                        showCloseButton: true,
                    });
                    clearInterval(intervalId);
                }
            }, 1000);
            let intrId = setInterval(() => {
                if (source === "primevideo") {
                    let x = document.querySelector(".swal2-container");
                    if (x) {
                        x.style.zIndex = "10000";
                        clearInterval(intrId);
                    }
                }
            }, 500);

        }
    }
}


function startListener(url) {
    console.log("%cMean-titles : Listener is working... ", "color: red;");
    let callback = () => {
        // let url = location.href;
        if (url.includes("netflix") || url.includes("youtube")) {
            const video = document.querySelector("video:first-of-type");
            // .player-timedText
            if (video && !netflixInit) {
                netflixInit = true;
                video.onpause = function (e) {
                    getMeaning( url.includes("netflix") ? "netflix": "youtube", video)
                }
            }
        } else if (url.includes("primevideo")) {
            let elem = document.querySelector('[id^="videoContainer_"]');
            if (elem) {
                const video = elem.firstChild;
                if (video && !primeInit) {
                    primeInit = true;
                    video.onpause = function (e) {
                        getMeaning("primevideo", video)
                    }
                }
            }
        }
    };



    const observer = new MutationObserver(callback);
    observer.observe(document.body, {
        subtree: true,
        attributes: false,
        childList: true
    });
}

String.prototype.removeStopWords = function() {
    var x;
    var y;
    var word;
    var stop_word;
    var regex_str;
    var regex;
    var cleansed_string = this.valueOf();
    var stop_words = ['a',
        'about',
        'above',
        'across',
        'after',
        'again',
        'against',
        'all',
        'almost',
        'alone',
        'along',
        'already',
        'also',
        'although',
        'always',
        'among',
        'an',
        'and',
        'another',
        'any',
        'anybody',
        'anyone',
        'anything',
        'anywhere',
        'are',
        'area',
        'areas',
        'around',
        'as',
        'ask',
        'asked',
        'asking',
        'asks',
        'at',
        'away',
        'b',
        'back',
        'backed',
        'backing',
        'backs',
        'be',
        'became',
        'because',
        'become',
        'becomes',
        'been',
        'before',
        'began',
        'behind',
        'being',
        'beings',
        'best',
        'better',
        'between',
        'big',
        'both',
        'but',
        'by',
        'c',
        'came',
        'can',
        'cannot',
        'case',
        'cases',
        'certain',
        'certainly',
        'clear',
        'clearly',
        'come',
        'could',
        'd',
        'did',
        'differ',
        'different',
        'differently',
        'do',
        'does',
        'done',
        'down',
        'down',
        'downed',
        'downing',
        'downs',
        'during',
        'e',
        'each',
        'early',
        'either',
        'end',
        'ended',
        'ending',
        'ends',
        'enough',
        'even',
        'evenly',
        'ever',
        'every',
        'everybody',
        'everyone',
        'everything',
        'everywhere',
        'f',
        'face',
        'faces',
        'fact',
        'facts',
        'far',
        'felt',
        'few',
        'find',
        'finds',
        'first',
        'for',
        'four',
        'from',
        'full',
        'fully',
        'further',
        'furthered',
        'furthering',
        'furthers',
        'g',
        'gave',
        'general',
        'generally',
        'get',
        'gets',
        'give',
        'given',
        'gives',
        'go',
        'going',
        'good',
        'goods',
        'got',
        'great',
        'greater',
        'greatest',
        'group',
        'grouped',
        'grouping',
        'groups',
        'h',
        'had',
        'has',
        'have',
        'having',
        'he',
        'her',
        'here',
        'herself',
        'high',
        'high',
        'high',
        'higher',
        'highest',
        'him',
        'himself',
        'his',
        'how',
        'however',
        'i',
        'if',
        'important',
        'in',
        'interest',
        'interested',
        'interesting',
        'interests',
        'into',
        'is',
        'it',
        'its',
        'itself',
        'j',
        'just',
        'k',
        'keep',
        'keeps',
        'kind',
        'knew',
        'know',
        'known',
        'knows',
        'l',
        'large',
        'largely',
        'last',
        'later',
        'latest',
        'least',
        'less',
        'let',
        'lets',
        'like',
        'likely',
        'long',
        'longer',
        'longest',
        'm',
        'made',
        'make',
        'making',
        'man',
        'many',
        'may',
        'me',
        'member',
        'members',
        'men',
        'might',
        'more',
        'most',
        'mostly',
        'mr',
        'mrs',
        'much',
        'must',
        'my',
        'myself',
        'n',
        'necessary',
        'need',
        'needed',
        'needing',
        'needs',
        'never',
        'new',
        'new',
        'newer',
        'newest',
        'next',
        'no',
        'nobody',
        'non',
        'noone',
        'not',
        'nothing',
        'now',
        'nowhere',
        'number',
        'numbers',
        'o',
        'of',
        'off',
        'often',
        'old',
        'older',
        'oldest',
        'on',
        'once',
        'one',
        'only',
        'open',
        'opened',
        'opening',
        'opens',
        'or',
        'order',
        'ordered',
        'ordering',
        'orders',
        'other',
        'others',
        'our',
        'out',
        'over',
        'p',
        'part',
        'parted',
        'parting',
        'parts',
        'per',
        'perhaps',
        'place',
        'places',
        'point',
        'pointed',
        'pointing',
        'points',
        'possible',
        'present',
        'presented',
        'presenting',
        'presents',
        'problem',
        'problems',
        'put',
        'puts',
        'q',
        'quite',
        'r',
        'rather',
        'really',
        'right',
        'right',
        'room',
        'rooms',
        's',
        'said',
        'same',
        'saw',
        'say',
        'says',
        'second',
        'seconds',
        'see',
        'seem',
        'seemed',
        'seeming',
        'seems',
        'sees',
        'several',
        'shall',
        'she',
        'should',
        'show',
        'showed',
        'showing',
        'shows',
        'side',
        'sides',
        'since',
        'small',
        'smaller',
        'smallest',
        'so',
        'some',
        'somebody',
        'someone',
        'something',
        'somewhere',
        'state',
        'states',
        'still',
        'still',
        'such',
        'sure',
        't',
        'take',
        'taken',
        'than',
        'that',
        'the',
        'their',
        'them',
        'then',
        'there',
        'therefore',
        'these',
        'they',
        'thing',
        'things',
        'think',
        'thinks',
        'this',
        'those',
        'though',
        'thought',
        'thoughts',
        'three',
        'through',
        'thus',
        'to',
        'today',
        'together',
        'too',
        'took',
        'toward',
        'turn',
        'turned',
        'turning',
        'turns',
        'two',
        'u',
        'under',
        'until',
        'up',
        'upon',
        'us',
        'use',
        'used',
        'uses',
        'v',
        'very',
        'w',
        'want',
        'wanted',
        'wanting',
        'wants',
        'was',
        'way',
        'ways',
        'we',
        'well',
        'wells',
        'went',
        'were',
        'what',
        'when',
        'where',
        'whether',
        'which',
        'while',
        'who',
        'whole',
        'whose',
        'why',
        'will',
        'with',
        'within',
        'without',
        'work',
        'worked',
        'working',
        'works',
        'would',
        'x',
        'y',
        'year',
        'years',
        'yet',
        'you',
        'young',
        'younger',
        'youngest',
        'your',
        'yours',
        'z']

    // Split out all the individual words in the phrase
    words = cleansed_string.match(/[^\s]+|\s+[^\s+]$/g)

    // Review all the words
    for(x=0; x < words.length; x++) {
        // For each word, check all the stop words
        for(y=0; y < stop_words.length; y++) {
            // Get the current word
            word = words[x].replace(/\s+|[^a-z]+/ig, "");	// Trim the word and remove non-alpha

            // Get the stop word
            stop_word = stop_words[y];

            // If the word matches the stop word, remove it from the keywords
            if(word.toLowerCase() === stop_word) {
                // Build the regex
                regex_str = "^\\s*"+stop_word+"\\s*$";		// Only word
                regex_str += "|^\\s*"+stop_word+"\\s+";		// First word
                regex_str += "|\\s+"+stop_word+"\\s*$";		// Last word
                regex_str += "|\\s+"+stop_word+"\\s+";		// Word somewhere in the middle
                regex = new RegExp(regex_str, "ig");

                // Remove the word from the keywords
                cleansed_string = cleansed_string.replace(regex, " ");
            }
        }
    }
    return cleansed_string.replace(/^\s+|\s+$/g, "");
};
