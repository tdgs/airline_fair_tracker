#!/usr/bin/env sh


function days_ahead {
  for d in $(seq 1 5)
  do
    date -d "now $d days" +%Y-%m-%d
  done
}

echo $(days_ahead 5) | xargs  -P5 -n 1 ./airline_fair_tracker -date
