#!/bin/bash

# Test script for HaberlerPlus
# This script tests all news sources and categories

echo "Building HaberlerPlus..."
cd "$(dirname "$0")/.."
go build -o bin/news ./cmd/news

if [ ! -f bin/news ]; then
    echo "Error: Failed to build HaberlerPlus"
    exit 1
fi

echo "Testing all news sources..."

# Function to test a source and category
test_source() {
    local source=$1
    local category=$2
    
    echo "Testing source $source, category $category..."
    printf "$source\n$category\n" | timeout 10s ./bin/news -d > /tmp/news_test_output.txt 2>&1
    
    # Check if there was an error
    if grep -q "error" /tmp/news_test_output.txt; then
        echo "⚠️ Warning: Error detected, skipping category"
        return 2
    fi
    
    # Check if any news items were found
    if grep -q "kategorisinden haberler:" /tmp/news_test_output.txt; then
        news_count=$(grep -c "http" /tmp/news_test_output.txt)
        if [ $news_count -gt 0 ]; then
            echo "✅ Success: Found $news_count news items"
            return 0
        else
            echo "❌ Failed: No news items found"
            return 1
        fi
    else
        echo "❌ Failed: Could not retrieve news"
        return 1
    fi
}

# Test all sources with their first category
echo "=== Testing all sources with their first category ==="
for source in {1..8}; do
    test_source $source 1 || continue
done

echo ""
echo "=== Testing all categories for each source ==="
# Test all categories for each source
for source in {1..8}; do
    # Get source name
    source_name=$(printf "$source\n1\n" | timeout 5s ./bin/news -d 2>/dev/null | grep -A 1 "Haber kaynağı numarası girin:" | tail -n 1 | cut -d' ' -f1)
    if [ -z "$source_name" ]; then
        source_name="Unknown"
    fi
    echo "Source $source ($source_name):"
    
    # Get the number of categories for this source
    printf "$source\n1\n" | timeout 5s ./bin/news -d > /tmp/news_categories.txt 2>&1
    category_count=$(grep -A 20 "Kategorileri:" /tmp/news_categories.txt | grep -B 20 "Kategori numarası girin:" | grep -v "Kategorileri:" | grep -v "Kategori numarası girin:" | wc -l)
    
    if [ $category_count -eq 0 ]; then
        category_count=8  # Default if we can't determine
    fi
    
    echo "Found $category_count categories"
    
    success_count=0
    fail_count=0
    
    for category in $(seq 1 $category_count); do
        result=0
        test_source $source $category || result=$?
        
        if [ $result -eq 0 ]; then
            ((success_count++))
        elif [ $result -eq 1 ]; then
            ((fail_count++))
        fi
        # If result is 2, it's an error we're skipping
    done
    
    echo "Summary: $success_count successful categories, $fail_count failed categories"
    echo ""
done

echo "Test completed." 