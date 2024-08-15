# moneybags

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

A Helm chart for Kubernetes

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| oci://registry-1.docker.io/bitnamicharts | postgresql | 15.5.21 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Pod Affinity |
| fullnameOverride | string | `""` | Full name override |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.repository | string | `"ghcr.io/ashwinath/moneybags"` | Repository |
| image.tag | string | `"latest"` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | Image pull secrets |
| moneybags.assets | string | `"date,type,amount\n2020-03-31,CPF,1000\n2020-03-31,Bank,20000\n2020-03-31,Mortgage,-40000\n2020-03-31,Investments,20000"` | CSV values for the assets |
| moneybags.expenses | string | `"date,type,amount\n2020-03-31,Credit Card,500\n2020-03-31,Reimbursement,-200\n2020-03-31,Tithe,800"` | CSV values for the expenses |
| moneybags.financials.alphavantageAPIKey | string | `"changeme"` | Alphavantage API key, get from https://www.alphavantage.co/support/#api-key |
| moneybags.financials.runIntervalInHours | int | `4` | Run financials job data population every x hours |
| moneybags.income | string | `"date,type,amount\n2021-03-11,Base,500\n2021-03-11,Bonus,200"` | CSV values for the income |
| moneybags.mortgage | string | `"mortgages:\n- total: 50000.0\n  mortgage_first_payment: 2022-10-10\n  mortgage_duration_in_years: 25\n  mortgage_end_date: 2047-10-10\n  interest_rate_percentage: 2.6\n  downpayments:\n  - date: 2021-10-10\n    sum: 1000.0\n  - date: 2021-12-12\n    sum: 20000.0"` | YAML values for mortgage |
| moneybags.shared_expenses | string | `"date,type,amount\n2023-01-01,Special:Renovations,5000.00\n2023-01-01,Electricity,100.00\n2023-01-01,Water,50.00\n2023-01-01,Gas,30.00\n2023-01-01,Grocery,300.00\n2023-01-01,Eating Out,500.00"` | CSV values for shared expenses |
| moneybags.telegram.allowedUser | string | `"changeme"` | Telegram user that will be allowed to send commands |
| moneybags.telegram.apiKey | string | `"changeme"` | Telegram bot API Key. Get by talking to Bot Godfather: https://core.telegram.org/bots/tutorial#obtain-your-bot-token |
| moneybags.telegram.debug | bool | `false` | Enable debugging. Warning: non standard logger used in telegram library |
| moneybags.trades | string | `"date_purchased,symbol,trade_type,price_each,quantity\n2021-03-11,IWDA.LON,buy,76.34,10"` | CSV values for the trades |
| nameOverride | string | `""` | Name override |
| nodeSelector | object | `{}` | Node Selectors |
| podAnnotations | object | `{}` | Pod Annotations |
| podLabels | object | `{}` | Pod Labels |
| podSecurityContext | object | `{}` | Pod Security Context |
| postgresql.auth.postgresPassword | string | `"changeme"` | Password for postgresql database, highly recommended to change this value |
| postgresql.primary.persistence.enabled | bool | `true` | Persist Postgresql data in a Persistent Volume Claim  |
| postgresql.resources | object | `{}` | Resources requests and limits for the database |
| replicaCount | int | `1` |  |
| resources | object | `{}` | Resources requests and limits for the moneybags |
| securityContext | object | `{}` | Security Context |
| tolerations | list | `[]` | Pod Tolerations |
