(function() {
    const protocol = 'http';
    const domain  = 'localhost:9996'
    const version = 'v1';
    let corpora = 'corpus';

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

    window.arb = window.arb || {};
    window.arb.q = window.arb.q || [];

    // initialization of window
    let elemDiv = document.createElement('div');
    elemDiv.style.cssText = 'width:20px;height:20px;background: red;position:absolute;bottom: 25px;right:25px;z-index:999999;border-radius: 30px;color: #fff;padding: 20px;font-weight: 800;font-family: "Arial";cursor: pointer;';
    elemDiv.textContent = "SP#";
    document.body.appendChild(elemDiv);

    // create window
    let elemWindow = document.createElement('div');
    elemWindow.style.cssText = 'display: none;width:200px;height:300px;background: #fff;position:absolute;bottom: 95px;right:25px;z-index:999999;border-radius: 20px;color: #000;padding: 20px;font-weight: 300;font-family: "Arial";border: 1px solid silver;overflow-x: auto;';
    document.body.appendChild(elemWindow);

    // add response ul
    let elemWrap = document.createElement('div');
    elemWrap.style.height = "90%";
    elemWrap.style.overflowY = "auto";
    let ul = document.createElement("ul");
    ul.style.cssText = "list-style: none;padding: 0";
    ul.setAttribute('id', 'chatWindow');
    elemWrap.appendChild(ul);
    elemWindow.appendChild(elemWrap);
    let elemBrand = document.createElement('div');
    elemBrand.textContent = "arualBot by ajandera.com";
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

    // input
    let elemInput = document.createElement("div");
    elemInput.innerHTML = "<input type='text' name='question' id='ask' value='' placeholder='How can I help you?' style='width: 95%;padding: 5px;position: absolute;left: 3%;bottom: 7px;border: none;border-top: 1px solid silver;'/>";
    elemWindow.appendChild(elemInput);

    // send btn
    let elemBtn = document.createElement("button");
    elemBtn.textContent = "Send"
    elemBtn.style.cssText = 'position: absolute;bottom: 2%;right: 5px;';

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
                let q = res.answer.replace("I am not able to anser your question :( Please teach how to answer on '", "").replace("'.", "");
                window.localStorage.setItem("arualbot_question", q);
            } else {
                window.localStorage.removeItem("arualbot_question");
            }
            let answer;
            if (res.success === true) {
                answer = document.createElement("li");
                answer.style.textAlign = "right";
                answer.style.fontStyle = "italic";
                answer.style.color = "silver";
                answer.style.paddingTop = "3px";
                answer.textContent = res.answer;
                document.getElementById('chatWindow').appendChild(answer);
            } else {
                answer = document.createElement("li");                
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

    elemBtn.onclick = () => {
        send();
    };
    elemWindow.appendChild(elemBtn);
    document.addEventListener("keyup", function(event) {
        if (event.keyCode === 13) {
            send();
        }
    });
}) ()