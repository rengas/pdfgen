package playground_test

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestTemplate(t *testing.T) {
	vals := map[string]interface{}{
		"amount":  10.2,
		"name":    "2312",
		"address": 12,
		"items":   []int{1, 2, 3},
		"itemMap": map[string]interface{}{"amount": 10.2},
	}
	dt := `<!DOCTYPE html>
<html>
   <head>
      <title>{{.amount}}</title>
   </head>
   <body>
      <h1>{{.amount}} </h1>
      <h1>{{.name}} </h1>
      <h1>{{.address}} </h1>
      <ul >
         {{range $i, $a := .items}}
         <li>{{$a}}</li>
         {{end}}
      </ul>
      <ul >
         {{range $i, $a := .itemMap}}
         <li>{{$a}}</li>
         {{end}}
      </ul>
   </body>
</html>`
	tr, err := template.New("test").Parse(string(dt))
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = tr.ExecuteTemplate(&b, "test", vals)
	if err != nil {
		panic(err)
	}

	fmt.Println(b.String())

}

func TestTemplate2(t *testing.T) {
	vals := map[string]interface{}{
		"logo":          "https://",
		"invoiceNumber": "INVOICE 3-2-1",
		"invoiceDetails": map[string]interface{}{
			"projectName": "Website development",
			"client":      "John Doe",
			"email":       "john@example.com",
			"date":        "August 17, 2015",
			"due date":    "September 17, 2015",
		},
		"invoiceAddress": map[string]interface{}{
			"companyName": "Company Name LLC",
			"streetName1": "455 Foggy Heights",
			"streetName2": "AZ 85004, US",
			"phoneNumber": "(602) 519-0450",
		},
		"lineItems": []map[string]interface{}{{
			"service":     "Design",
			"description": "Creating a recognizable design solution based on the company's existing visual identity",
			"unitPrice":   "$40.00",
			"quantity":    26,
			"total":       "$1,040.00",
		},
			{
				"service":     "Development",
				"description": "Developing a Content Management System-based Website",
				"unitPrice":   "$40.00",
				"quantity":    80,
				"total":       "$3,200.00",
			},
			{
				"service":     "SEO",
				"description": "Optimize the site for search engines (SEO)",
				"unitPrice":   "$40.00",
				"quantity":    20,
				"total":       "$800.00",
			},
			{
				"service":     "Training",
				"description": "Initial training sessions for staff responsible for uploading web content",
				"unitPrice":   "$40.00",
				"quantity":    4,
				"total":       "$160.00",
			},
		},
		"subtotal":   "$5,200.00",
		"tax":        "$1,300.00",
		"grandTotal": "$6,500.00",
	}

	dt := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Example 1</title>
    <style>
        .clearfix:after {
            content: "";
            display: table;
            clear: both;
        }
        a {
            color: #5D6975;
            text-decoration: underline;
        }
        body {
            position: relative;
            width: 21cm;
            height: 29.7cm;
            margin: 0 auto;
            color: #001028;
            background: #FFFFFF;
            font-family: Arial, sans-serif;
            font-size: 12px;
            font-family: Arial;
        }
        header {
            padding: 10px 0;
            margin-bottom: 30px;
        }
        #logo {
            text-align: center;
            margin-bottom: 10px;
        }
        #logo img {
            width: 90px;
        }
        h1 {
            border-top: 1px solid  #5D6975;
            border-bottom: 1px solid  #5D6975;
            color: #5D6975;
            font-size: 2.4em;
            line-height: 1.4em;
            font-weight: normal;
            text-align: center;
            margin: 0 0 20px 0;
            background: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACsAAAAyCAMAAADyWtKhAAAAOVBMVEXu7u7////v7+/x8fHw8PD+/v79/f38/Pz29vb39/fy8vL09PTt7e319fX6+vr5+fn4+Pjz8/P7+/sjQsgnAAABrklEQVR4Xu2Wy2pkMQxET0ny+z77/z92EHRIM4RJejeLeGGMJbkKu1QYt0vVbajkJJWjTw+f/SiSBl40zKsuc5hVDzjVjEtlOB/DR9GFNZ3wUJ1wFmkRVQe0NoFYffQVwGwNDtVgSeUkgdh2acePAM9iZcghDidDG0mPetzgVdXwgJ6ZpeRcOoRjVdXhPipu4E1aYHBJOpeHr1PSBUbiNwdzDNtrbhvQpXbDtm1wN6kDxiXV3XJx1DzBCGMWtQDfx+4QTWVigSVaPS4kqS2IbpwqNySM6g530Yn1gNUkidLOFeCtExmjSypK/KwNenOwdbaCG7msurclLbyorlhVxcmN7VZdBphjZr4/pBJbz/iucm9sd9GedX2LIj12NzNi9uSXueNZEFjyGdvHlPz7jHfO/eQ7+ZIv85Nv3oP97B7eud/nu41/vtt4vttTDwPbYPylhwGbMZ56eNUZX+iMF5296nd+od/5ot93+uKtfrOXPj5VRvAxYhSdL31sBFcWDqUgEvvRZ6RKHkV6aiJhL+I1l+9y7cdeYv+pl/x6ya+X/HoJP/YSXv8afPPX+AONyic1BlYxVgAAAABJRU5ErkJggg==');
        }
        #project {
            float: left;
        }
        #project span {
            color: #5D6975;
            text-align: right;
            width: 52px;
            margin-right: 10px;
            display: inline-block;
            font-size: 0.8em;
        }
        #company {
            float: right;
            text-align: right;
        }

        #project div,
        #company div {
            white-space: nowrap;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            border-spacing: 0;
            margin-bottom: 20px;
        }

        table tr:nth-child(2n-1) td {
            background: #F5F5F5;
        }

        table th,
        table td {
            text-align: center;
        }

        table th {
            padding: 5px 20px;
            color: #5D6975;
            border-bottom: 1px solid #C1CED9;
            white-space: nowrap;
            font-weight: normal;
        }

        table .service,
        table .desc {
            text-align: left;
        }

        table td {
            padding: 20px;
            text-align: right;
        }

        table td.service,
        table td.desc {
            vertical-align: top;
        }

        table td.unit,
        table td.qty,
        table td.total {
            font-size: 1.2em;
        }

        table td.grand {
            border-top: 1px solid #5D6975;;
        }

        #notices .notice {
            color: #5D6975;
            font-size: 1.2em;
        }

        footer {
            color: #5D6975;
            width: 100%;
            height: 30px;
            position: absolute;
            bottom: 0;
            border-top: 1px solid #C1CED9;
            padding: 8px 0;
            text-align: center;
        }
    </style>
</head>
<body>
<header class="clearfix">
    <div id="logo">
        <img src="logo.png">
    </div>
    <h1>{{.invoiceNumber}}</h1>
    <div id="company" class="clearfix">
        <div>{{ .invoiceAddress.companyName}}</div>
        <div>{{ .invoiceAddress.streetName1}}<br /> {{ .invoiceAddress.streetName2}}</div>
        <div>{{ .invoiceAddress.phoneNumber}}</div>
        <div><a href="mailto:{{ .invoiceAddress.email}}"> {{ .invoiceAddress.email}}</a></div>
    </div>
    <div id="project">
        <div><span>PROJECT</span> {{ .invoiceDetails.projectName}}</div>
        <div><span>CLIENT</span> {{ .invoiceDetails.client}}</div>
        <div><span>ADDRESS</span> {{ .invoiceDetails.address}}</div>
        <div><span>EMAIL</span> <a href="mailto:{{ .invoiceDetails.email}}"> {{ .invoiceDetails.email}}</a><</div>
        <div><span>DATE</span> {{ .invoiceDetails.date}}</div>
        <div><span>DUE DATE</span> {{ .invoiceDetails.dueDate}}</div>
    </div>
</header>
<main>
    <table>
        <thead>
        <tr>
            <th class="service">SERVICE</th>
            <th class="desc">DESCRIPTION</th>
            <th>PRICE</th>
            <th>QTY</th>
            <th>TOTAL</th>
        </tr>
        </thead>
        <tbody>
        {{range $i, $item := .lineItems}}
        <tr>
            <td class="service">{{$item.service}}</td>
            <td class="desc">{{$item.description}}</td>
            <td class="unit">{{$item.unitPrice}}</td>
            <td class="qty">{{$item.quantity}}</td>
            <td class="total">{{$item.total}}</td>
        </tr>
        {{end}}
        <tr>
            <td colspan="4">SUBTOTAL</td>
            <td class="total">{{.subTotal}}</td>
        </tr>
        <tr>
            <td colspan="4">TAX 25%</td>
            <td class="total">{{.tax}}</td>
        </tr>
        <tr>
            <td colspan="4" class="grand total">GRAND TOTAL</td>
            <td class="grand total">{{.grandTotal}}</td>
        </tr>
        </tbody>
    </table>
    <div id="notices">
        <div>NOTICE:</div>
        <div class="notice">A finance charge of 1.5% will be made on unpaid balances after 30 days.</div>
    </div>
</main>
<footer>
    Invoice was created on a computer and is valid without the signature and seal.
</footer>
</body>
</html>`
	tr, err := template.New("test").Parse(string(dt))
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = tr.ExecuteTemplate(&b, "test", vals)
	if err != nil {
		panic(err)
	}

	fmt.Println(b.String())

}
