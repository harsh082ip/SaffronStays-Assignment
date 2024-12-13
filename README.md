# SaffronStays-Assignment

## Objective

A REST API for fetching Airbnb room performance metrics, including occupancy percentage and night rates for the next 5 months and 30 days.

## Requirements

- Framework: Go-Gin 
- Database: PostgreSQL

## Running the Project

### Docker Deployment

```bash
docker run -e SERVICE_URI="postgres://<username>:<password>@<host>:<port>/<database>?sslmode=require" -p 8000:8000 harsh082ip/go-metrics-service
```

### Local Setup

1. Clone the repository:
```bash
git clone https://github.com/your-repo/go-metrics-service.git
cd go-metrics-service
```

2. Install dependencies and run:
```bash
go mod tidy
export SERVICE_URI="your-connection-string"
go run cmd/main.go
```

## API Endpoint

`GET /:room_id`

### Example Response
```json
{
  "room_id": 101,
  "metrics": {
    "night_rates": {
      "average_rate": 150.5,
      "highest_rate": 200.0,
      "lowest_rate": 120.0
    },
    "occupancy": [
      {
        "month": "January 2024",
        "occupancy_percentage": 64.52
      }
    ]
  }
}
```

## Technical Challenges and Solutions

### 1. Database Data Insertion

**Initial Approach**: 
- Originally, the project used individual database writes to insert 1500 room records.
- This approach was highly inefficient, resulting in:
  - Slow data population
  - Increased database connection overhead
  - Significant time consumption for inserting records

**Optimized Solution**: 
- Implemented bulk upload mechanism
- Reduced insertion time dramatically
- Minimized database connection overhead
- Improved overall data population performance

### 2. Date Handling Complexity

**Challenge**: Managing complex date parsing and storage

**Implementation Details**:
- Added `AvailableDatesRaw` field to the `Room` struct
- Purpose of `AvailableDatesRaw`:
  - Temporarily store raw date strings from the database
  - Provide an intermediate step for date conversion
  - Enable flexible parsing of date formats
  - Handle potential variations in date string representations

**Parsing Process**:
```go
// Convert raw date strings to time.Time objects
for _, dateStr := range room.AvailableDatesRaw {
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        // Handle parsing errors
    }
    room.AvailableDates = append(room.AvailableDates, date)
}
```

### 3. Occupancy Calculation

**Complex Calculation**:
- Iterate through available dates
- Group dates by month
- Calculate available days per month
- Compute occupancy percentage dynamically

**Key Techniques**:
- Used map for efficient month-based tracking
- Dynamically calculated total days in each month
- Precise percentage calculation with two decimal precision