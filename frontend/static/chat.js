(function() {
    window.addEventListener('DOMContentLoaded', () => {
    const protocol = 'http',
    domain  = 'localhost:3000',
    version = 'v1',
    css = '@keyframes blink {0%, 50% {opacity: 1;}50.1%, 100% {opacity: 0;}}';

    let corpora = 'corpus',
    elemDiv = document.createElement('div'),
    elemBtn = document.createElement("button"),
    elemWindow = document.createElement('div'),
    elemWrap = document.createElement('div'),
    ul = document.createElement("ul"),
    elemBrand = document.createElement('div'),
    elemInput = document.createElement("div"),
    answer = null,
    response = '',
    counter = 0,
    index = 0,
    typewriter =  null;

    window.arb = window.arb || {};
    window.arb.q = window.arb.q || [];

    // initialization of window
    elemDiv.style.cssText = 'background: #B1EA4E;position:fixed;bottom: 25px;right:25px;z-index:999999;border-radius: 30px;color: #4633af;padding: 20px;font-weight: 800;font-family: "Arial";cursor: pointer;';
    elemDiv.textContent = "|>";
    document.body.appendChild(elemDiv);

    // add css style
    head = document.head || document.getElementsByTagName('head')[0],
    style = document.createElement('style');
    head.appendChild(style);
    style.type = 'text/css';

    if (style.styleSheet) {
        // This is required for IE8 and below.
        style.styleSheet.cssText = css;
    } else {
        style.appendChild(document.createTextNode(css));
    }

    // create window
    elemWindow.style.cssText = 'display: none;width:500px;height:500px;background: #fff;position:fixed;bottom: 95px;right:25px;z-index:999999;border-radius: 20px;color: #000;padding: 20px;font-weight: 300;font-family: "Arial";border: 1px solid silver;overflow-x: auto;';
    document.body.appendChild(elemWindow);

    // add response ul
    elemWrap.style.height = "95%";
    elemWrap.style.overflowY = "auto";
    ul.style.cssText = "list-style: none;padding: 0";
    ul.setAttribute('id', 'chatWindow');
    elemWrap.appendChild(ul);
    elemWindow.appendChild(elemWrap);
    elemBrand.textContent = "powered by storePredictor AI";
    elemBrand.style.cssText = "position: absolute;top: 2px;right: 10px;font-size: 10px;";
    elemWindow.appendChild(elemBrand);

    // toggle chat
    elemDiv.onclick = () => {
        if (elemWindow.style.display === 'block' || elemWindow.style.display === '') {
            elemWindow.style.display = 'none';
        } else {
            elemWindow.style.display = 'block'
        }
    };

    document.addEventListener('click', event => {
        const isClickInside = elemDiv.contains(event.target)

        if (!isClickInside) {
            elemWindow.style.display = 'none';
        }
    })

    // input
    elemInput.style.width = "100%";
    elemInput.style.borderTop = "1px solid silver";
    elemInput.innerHTML = "<input type='text' name='question' id='ask' value='' placeholder='How can I help you?' style='width: 80%;padding: 5px;position: absolute;left: 3%;bottom: 7px;border: none;'/>";
    elemWindow.appendChild(elemInput);    

    // send btn
    elemBtn.textContent = "Send"
    elemBtn.style.cssText = 'position: absolute;bottom: 1%;right: 15px;background: #00c9db;color: #fff;padding: 3px 10px;border-radius: 6px;';
    elemBtn.onclick = () => {
        send();
    };

    elemWindow.appendChild(elemBtn);
    document.addEventListener("keyup", function(event) {
        if (event.keyCode === 13) {
            send();
        }
    });

    // send function
    function send() {
        let question = document.getElementById('ask').value;
        let li = document.createElement("li");
        li.style.paddingTop = "3px";
        li.textContent = question;
        document.getElementById('chatWindow').appendChild(li);
        document.getElementById('ask').value = "";
        document.getElementById('chatWindow').scrollIntoView({ behavior: 'smooth', block: 'end' });
        const questionTolearn = window.localStorage.getItem("arualbot_question");
        post("/ask", {
            'text': question,
            'isTeach': questionTolearn !== null,
            'question': questionTolearn !== null ? questionTolearn : "",
            'corpora': corpora,
        }).then((res) => {
            // check if need to learn
            if (res.needLearn === true) {
                let q = res.answer.replace("I am not able to anser your question :( Please teach me how to answer on '", "").replace("'.", "");
                window.localStorage.setItem("arualbot_question", q);
            } else {
                window.localStorage.removeItem("arualbot_question");
            }
            
            if (res.success === true) {
                counter++;
                answer = document.createElement("li");
                answer.setAttribute("id", 'typewriter'+counter);
                //answer.style.textAlign = "right";
                answer.style.fontStyle = "italic";
                answer.style.color = "gray";
                answer.style.paddingTop = "5px";
                answer.style.paddingBottom = "5px";
                response = parseResponse(res.answer);
                document.getElementById('chatWindow').appendChild(answer);
                typewriter = document.getElementById('typewriter'+counter);
                index = 0;
                // type only if response not contain html
                if (resnponseText.type) {
                    type();
                } else {
                    setTimeout(() => {
                        answer.innerHTML = response;
                    }, 2000).then(() => document.getElementById('chatWindow').scrollIntoView({ behavior: 'smooth', block: 'end' }));
                }
            } else {
                answer = document.createElement("li");
                answer.className = 'typewriter';       
                answer.style.textAlign = "right";
                answer.style.fontStyle = "italic";
                answer.style.paddingTop = "3px";
                answer.style.color = "red";
                answer.textContent = "An error occured."
                document.getElementById('chatWindow').appendChild(answer);
            }
            document.getElementById('chatWindow').scrollIntoView({ behavior: 'smooth', block: 'end' });
        });
    }

    // human type simulation
    function type() {
        if (index < response.length) {
            typewriter.innerHTML = response.slice(0, index) + '<span style="margin-left: 5px;background-color: #fff;animation: blink 1s infinite;">|</span>';
            index++;
            document.getElementById('chatWindow').scrollIntoView({ behavior: 'smooth', block: 'end' });
            setTimeout(type, Math.random() * 15 + 50);
        } else {
            typewriter.innerHTML = response.slice(0, index) + '<span style="margin-left: 5px;background-color: #fff;animation: blink 1s infinite;>|</span>';
        }
    }

    function parseResponse(text) {
        // detect link
        let urlRegex = /(https?:\/\/[^\s]+)/g;
        // detect img
        let imgRegex = /[\w-]+\.(jpg|png|txt)/g;
        const img = text.match(imgRegex);
        if (img !== null) {
            return text.replace(img[0], '<br><img src="/'+img[0]+'" style="width:450px;cursor:pointer">').replace(urlRegex, '<a href="$1">$1</a>');
        } else {
            return text.replace(urlRegex, '<a href="$1">$1</a>');
        }
    }

    async function post(endpoint, data) {
        let a = protocol + "://" + domain + "/" + version + endpoint;
        //data['code'] = '';

        if (typeof (fetch) !== 'function') {
            return null;
        }

        const response = await fetch(a, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        const jsonData = await response.json();
        return jsonData;
    };
})
}) ()