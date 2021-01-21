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
                    sendCustomerDialog("updateCustomer");
                },
            },
        });
    });

    //кнопка клиента Удалить
    $('#btc_delete').on('click', function() {
        customerDeleteB();
    });


    //кнопка лицензия Создание
    $('#btl_create').on('click', function() {
        setCreateLicenseDialog();
        $('#licDialog').dialog('open');
        $('#licDialog').dialog({
            buttons: {
                "Отправить": function() {
                    sendLicenseDialog("createLicense");
                },
            },
        });
    });

    //кнопка лицензии Обновить
    $('#btl_update').on('click', function() {
        setlicenseUpdateDialog();
        $('#licDialog').dialog('open');
        $('#licDialog').dialog({
            buttons: {
                "Отправить": function() {
                    sendLicenseDialog("updateLicense");
                },
            },
        });
    });

    //кнопка лицензия Удаление
    $('#btl_delete').on('click', function() {
        licenseDeleteB();
    });

    $('#btt_copy').on('click', function() {
        let selected = $('#tableLicense').bootstrapTable('getSelections');
        copyTextToBuffer(selected[0].token);
        successAlertMessage("Ключ скопирован");
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
    let noCheck = false;
    customers.forEach(cust => {
        let temp = {
            id: cust.id,
            check: false,
            name: cust.name,
            address: cust.address,
            numS: cust.licenses.length,
            phone: cust.phone,
            email: cust.email,
        };
        if (selected.length === 1) {
            if (cust.id === selected[0].id) {
                temp.check = true;
                noCheck = true;
                fillLicenseTalbe();
            }
        }
        toWrite.push(temp);
    });
    if (!noCheck) {
        setClientDisableBut(true);
        setLicenseDisableBut(true);
        $('#cName').text("");
    }
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

//setCreateLicenseDialog диалог при создании лицензии
function setCreateLicenseDialog() {
    $('#licDialog').dialog({
        autoOpen: false,
        resizable: false,
    });
    $('#numDev').val(1);
    $('#numAcc').val(1);
    $('#yakey').val("");
    $('#emailList').val("");
    $('#endTime').val(new Date().toISOString().slice(0, 10));
};

//setCustomerUpdateDialog диалог при обновлении клиента
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
};

//setlicenseUpdateDialog диалог при обновлении лицензии
function setlicenseUpdateDialog() {
    $('#licDialog').dialog({
        autoOpen: false,
        resizable: false,
    });
    let selected = $('#tableLicense').bootstrapTable('getSelections');
    $('#numDev').val(selected[0].numDev);
    $('#numAcc').val(selected[0].numAcc);
    $('#yakey').val(selected[0].yaKey);
    let time = selected[0].endTime.slice(0, 10).split('.');
    $('#endTime').val(time[2] + "-" + time[1] + "-" + time[0]);
    $('#emailList').val(selected[0].email.toString().replaceAll(',', ' '));
};

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

//sendLicenseDialog отпрака информации из диалога клиентов на сервер
function sendLicenseDialog(typeD) {
    let licForm = $('#licForm')
    if (!licForm[0].checkValidity()) {
        licForm[0].classList.add('was-validated');
        return
    }
    let selectCust = $('#table').bootstrapTable('getSelections');
    let selectLic = $('#tableLicense').bootstrapTable('getSelections');

    let toSend = {
        type: typeD,
        idCust: selectCust[0].id,
        license: {
            numdev: parseInt($('#numDev').val()),
            numacc: parseInt($('#numAcc').val()),
            yakey: $('#yakey').val(),
            endtime: (new Date($('#endTime').val())).toISOString(),
            tech_email: $('#emailList').val().split(" "),
        },
    };

    if (selectLic.length === 0) {
        toSend.license.id = 0;
    } else {
        toSend.license.id = selectLic[0].id;
    }
    console.log(toSend);
    ws.send(JSON.stringify(toSend));
    $('#licDialog').dialog('close');
};


//deleteB удаление клиента
function customerDeleteB() {
    let selected = $('#table').bootstrapTable('getSelections');
    let toSend = {
        type: "deleteCustomer",
        id: selected[0].id,
    };
    setClientDisableBut(true);
    ws.send(JSON.stringify(toSend));
};

//deleteB удаление лицензии
function licenseDeleteB() {
    let selCust = $('#table').bootstrapTable('getSelections');
    let selLic = $('#tableLicense').bootstrapTable('getSelections');
    let toSend = {
        type: "deleteLicense",
        idCust: selCust[0].id,
        license: {
            id: selLic[0].id,
        },
    };
    setLicenseDisableBut(true);
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
            yaKey: lic.yakey,
            token: lic.token,
            endTime: timeFormat(lic.endtime),
        };
        console.log(timeFormat(lic.endtime));
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
        year: "numeric",
        hour: "numeric",
        minute: "numeric",
        second: "numeric",
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