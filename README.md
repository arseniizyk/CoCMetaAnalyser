# CoCMetaAnalyser

A tool that retrieves Hero Equipment data from the top 10,000 Clash of Clans players.

## [Visual Charts](https://arseniizyk.github.io/CoCMetaAnalyser/)

## 📂 Data Storage

You can find the latest collected data in the [`cmd/CoCMetaAnalyser/meta`](cmd/CoCMetaAnalyser/meta) directory.  
The data is updated monthly, starting from January 2025.

## 🚀 How to Run

1. **Set Up API Access**

    - Register at [Clash of Clans API](https://developer.clashofclans.com/)
    - Get your API token

2. **Create a `.env` File**

    - Move it in `cmd/CoCMetaAnalyser/`
    - Add fields from [`.env.example`](.env.example):
        ```ini
        COC_API=your_token
        ```

3. **Run**  
   Execute the application using `go run main.go` and use meta or metapairs flags
    ```bash
     go run main.go meta -season 2025-05 -limit 100 -filename meta
    ```

## 📌 Notes

-   Ensure your API token is valid before running the tool.
-   Data updates are performed manually every month.
-   You can increase or decrease the number of players (minimum: 100, maximum: 25,000) in [`app.go`](internal/app/app.go).
-   I have also created some functions that are not currently used, but the API provides them (retrieving all available seasons and league information).

## 📁 Data

-   ## 2025
    -   January [Items](cmd/app/meta/2025-01)
    -   February [Items](cmd/app/meta/2025-02) | [Pairs](cmd/app/metapairs/2025-02)
    -   March [Items](cmd/app/meta/2025-03) | [Pairs](cmd/app/metapairs/2025-03)
    -   April [Items](cmd/app/meta/2025-04) | [Pairs](cmd/app/metapairs/2025-04)
    -   May [Items](cmd/app/meta/2025-05) | [Pairs](cmd/app/metapairs/2025-05)
    -   June [Items](cmd/app/meta/2025-06) | [Pairs](cmd/app/metapairs/2025-06)
    -   July [Items](cmd/app/meta/2025-07) | [Pairs](cmd/app/metapairs/2025-07)

---

## 🔧 Planned Features

-   ✅ Improve performance by using Goroutines and concurrency for fetching player data(2025-02-09)
-   ✅ CLI usage(2025-05-07)
