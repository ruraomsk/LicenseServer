'use strict';
let ws;
let customers;

$(document).ready(function() {

    let url = 'ws://' + location.host + location.pathname + location.search + 'custMain/ws'
    ws = new WebSocket(url)

    ws.onopen = function() {
        console.log('connected')
    };

    ws.onclose = function(evt) {
        console.log('disconnected', evt);
    };

    ws.onmessage = function(evt) {
        let allData = JSON.parse(evt.data);
        let data = allData.data;
        console.log(allData);

        switch (allData.type) {
            case 'custInfo':
                {
                    customers = data.custInfo.sort(sortByName);
                    fillCustomerTalbe();
                    break;
                }
            case 'custUpdate':
                {
                    customers = data.custInfo.sort(sortByName);
                    fillCustomerTalbe();
                    break;
                }
            case 'error':
                {
                    $('#warningAlert').text(data.message.error)
                    $('#warningAlert').fadeTo(3000, 500).slideUp(500, function() {
                        $('#warningAlert').slideUp(500);
                    });
                    break;
                }

            default:
                break;
        };
    };


    setClientDisableBut(true);
    setLicenseDisableBut(true);

    //кнопка клиента Создать
    $('#btc_create').on('click', function() {
        setCreateCustomerDialog();
        $('#custDialog').dialog('open');
        $('#custDialog').dialog({
            buttons: {
                "Отправить": function() {
                    sendCustomerDialog("createCustomer");
                },
            },
        });
    });

    //кнопка клиента Обновить
    $('#btc_update').on('click', function() {
        setCustomerUpdateDialog();
        $('#custDialog').dialog('open');
        $('#custDialog').dialog({
            buttons: {
                "Отправить": function() {
                    sendCustomerDialog("updateCustomer")
                },
            },
        });
    });

    //кнопка клиента Удалить
    $('#btc_delete').on('click', function() {
        customerDeleteB();
    });

    $('#btt_copy').on('click', function() {
        let selected = $('#tableLicense').bootstrapTable('getSelections');
        copyTextToBuffer(selected[0].token);
        successAlertMessage("Сообщение скопированно");
    });
});

function successAlertMessage(message) {
    $('#successAlert').text(message);
    $('#successAlert').fadeTo(3000, 500).slideUp(500, function() {
        $('#successAlert').slideUp(500);
    });
}

function sortByName(a, b) {
    return a.name - b.name;
};

//fillCustomerTalbe заполнить таблицу клиентов
function fillCustomerTalbe() {
    let $table = $('#table');
    let selected = $table.bootstrapTable('getSelections');
    let toWrite = [];
    customers.forEach(cust => {
        let temp = {
            id: cust.id,
            check: false,
            name: cust.name,
            address: cust.address,
            numS: cust.servers.length,
            phone: cust.phone,
            email: cust.email,
        };
        if (selected.length === 1) {
            if (cust.id === selected[0].id) {
                temp.check = true;
            }
        }
        toWrite.push(temp);
    });
    $table.bootstrapTable('load', toWrite);
    $table.bootstrapTable('hideColumn', 'id');
    $table.bootstrapTable('scrollTo', 'top');

    $table.on('click', function() {
        let selected = $table.bootstrapTable('getSelections');
        if (selected.length > 0) {
            setClientDisableBut(false);
            $('#cName').text(selected[0].name);
            fillLicenseTalbe();
        } else {
            setClientDisableBut(true);
            setLicenseDisableBut(true);
            $('#cName').text("");
        }
    });
};

//setDisableButtons устанавливает неактивность кнопок
function setClientDisableBut(flag) {
    $('#btc_update').prop('disabled', flag); //клиент обновление
    $('#btc_delete').prop('disabled', flag); //клиент удаление
    $('#btl_create').prop('disabled', flag); //лицензия создание
    $('#tableLicense').bootstrapTable('load', []);
};

//setCreateDialog диалог при создании клиента
function setCreateCustomerDialog() {
    $('#custDialog').dialog({
        autoOpen: false,
        resizable: false,
    });
    $('#name').val("");
    $('#address').val("");
    $('#phone').val("");
    $('#email').val("");
};

//setUpdateDialog диалог при обновлении клиента
function setCustomerUpdateDialog() {
    $('#custDialog').dialog({
        autoOpen: false,
        resizable: false,
    });
    let selected = $('#table').bootstrapTable('getSelections');
    $('#name').val(selected[0].name);
    $('#address').val(selected[0].address);
    $('#phone').val(selected[0].phone);
    $('#email').val(selected[0].email);
}

//sendCustomerDialog отпрака информации из диалога клиентов на сервер
function sendCustomerDialog(typeD) {
    let cForm = $('#custForm')
    if (!cForm[0].checkValidity()) {
        cForm[0].classList.add('was-validated');
        return
    }
    let selected = $('#table').bootstrapTable('getSelections');
    let toSend = {
        type: typeD,
        name: $('#name').val(),
        address: $('#address').val(),
        phone: $('#phone').val(),
        email: $('#email').val(),
    };
    if (selected.length === 0) {
        toSend.id = 0;
    } else {
        toSend.id = selected[0].id;
    }
    ws.send(JSON.stringify(toSend));
    $('#custDialog').dialog('close');
};

//deleteB удаление клиента
function customerDeleteB() {
    let selected = $('#table').bootstrapTable('getSelections');
    let toSend = {
        type: "deleteCustomer",
        id: selected[0].id,
    };
    $('#btc_update').prop('disabled', true);
    $('#btc_delete').prop('disabled', true);
    ws.send(JSON.stringify(toSend));
};


//fillLicenseTalbe заполнить таблицу клиентов
function fillLicenseTalbe() {
    let selectedCust = $('#table').bootstrapTable('getSelections');
    let $table = $('#tableLicense');
    let selectedLic = $table.bootstrapTable('getSelections');
    let toWrite = [];
    customers.find(item => item.id === selectedCust[0].id).licenses.forEach(lic => {
        let temp = {
            id: lic.id,
            check: false,
            numDev: lic.numdev,
            numAcc: lic.numacc,
            email: lic.tech_email,
            token: lic.token,
            endTime: timeFormat(lic.endtime),
        };
        if (selectedLic.length === 1) {
            if (lic.id === selectedLic[0].id) {
                temp.check = true;
            }
        }
        toWrite.push(temp);
    });


    $table.bootstrapTable('load', toWrite);
    $table.bootstrapTable('hideColumn', 'id');
    $table.bootstrapTable('scrollTo', 'top');


    $table.on('click', function() {
        let selected = $table.bootstrapTable('getSelections');
        if (selected.length > 0) {
            setLicenseDisableBut(false);
        } else {
            setLicenseDisableBut(true);
        }
    });
};

//timeFormat преобразование переданного формата для отображения
function timeFormat(time) {
    let date = new Date(time);
    const dateTimeFormat = new Intl.DateTimeFormat('ru', {
        day: "2-digit",
        month: "2-digit",
        year: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        timeZoneName: "short"
    });
    return dateTimeFormat.format(date);
};


//setLicenseDisableBut установка отображения кнопок
function setLicenseDisableBut(flag) {
    $('#btl_update').prop('disabled', flag); //лицензия обновление
    $('#btl_delete').prop('disabled', flag); //лицензия удаление
    $('#btt_recreate').prop('disabled', flag); //ключ пересоздание
    $('#btt_copy').prop('disabled', flag); //ключ копирование
};

//copyTextToBuffer копирование текста в буфер обмена
function copyTextToBuffer(value) {
    var $temp = $("<input>");
    $("body").append($temp);
    $temp.val(value).select();
    document.execCommand("copy");
    $temp.remove();
};