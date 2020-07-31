'use strict';
let globalData;
let ws;


$(document).ready(function() {

    let url = 'ws://' + location.host + location.pathname + location.search + '/ws'
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
                    let customers = data.custInfo
                    fillTalbeCustomer(customers, true)
                    break;
                }
            case 'custUpdate':
                break;
            default:
                break;
        };
    };

    setCreateDialog();

    $('#bt_create').on('click', function() {
        $('#createDialog').dialog('open');
    });

    $('#bt_delete').on('click', function() {

    });
});


function fillTalbeCustomer(customers, firstFlag) {
    let $table = $('#table');
    let selected = $table.bootstrapTable('getSelections');
    let toWrite = [];
    customers.forEach(cust => {
        let temp = {
            id: cust.id,
            check: (selected.length !== 0) ? (cust.name === customers[0].name) : false,
            name: cust.name,
            address: cust.address,
            numS: cust.servers.length,
            phone: cust.phone,
            email: cust.email,
        };
        if ((cust.name === customers[0].name) && firstFlag) {
            temp.check = true;
        }
        toWrite.push(temp);
    });
    $table.bootstrapTable('load', toWrite)
    $table.bootstrapTable('hideColumn', 'id');
    $table.bootstrapTable('scrollTo', 'top');
};


function setCreateDialog() {
    $('#createDialog').dialog({
        autoOpen: false,
        buttons: {
            'Отправить': sendBCreateD,
            'Закрыть': function() {
                $(this).dialog('close');
            }
        },
        resizable: false,
    });

};



function sendBCreateD() {
    let cForm = $('#createForm')
    if (!cForm[0].checkValidity()) {
        cForm[0].classList.add('was-validated');
        return
    }

    let toSend = {
        type: "createCustomer",
        name: $('#name').val(),
        address: $('#address').val(),
        phone: $('#phone').val(),
        email: $('#email').val(),
    };
    ws.send(JSON.stringify(toSend));
    $(this).dialog('close');
}

function deleteB() {
    let selected = $('#table').bootstrapTable('getSelections');
    let id;
    selected
    customers.forEach(cust => {

    });

    let toSend = {
        type: "deleteCustomer",
        id: $().val(),
    }


}