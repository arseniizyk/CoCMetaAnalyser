# CoCMetaAnalyser

A tool that retrieves Hero Equipment data from the top 10,000 Clash of Clans players.

## 📂 Data Storage

You can find the latest collected data in the [`cmd/CoCMetaAnalyser/meta`](cmd/CoCMetaAnalyser/meta) directory.  
The data is updated monthly, starting from January 2025.

## 🚀 How to Run

1. **Set Up API Access**

    - Register at [Clash of Clans API](https://developer.clashofclans.com/)
    - Get your API token

2. **Create a `.env` File**

    - Place it in `cmd/CoCMetaAnalyser/`
    - Add fields from `.env.example`:
        ```ini
        COC_API=your_token
        ```

3. **Run**  
   Execute the application using `go run main.go`

## 📌 Notes

-   Ensure your API token is valid before running the tool.
-   Data updates are performed manually every month.
-   You can increase or decrease the number of players (minimum: 100, maximum: 25,000) in app.go.
-   I have also created some functions that are not currently used, but the API provides them (retrieving all available seasons and league information).

---

## 🔧 Planned Features  
- ✅ Improve performance by using Goroutines and concurrency for fetching player data(2025-02-09)
