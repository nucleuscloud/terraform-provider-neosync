data "neosync_job" "my_job" {
  id = "3b83d1d3-5ffe-48c6-ac11-7a2e60802864"
}

output "job_name" {
  value = data.neosync_job.my_job.name
}

output "job_id" {
  value = data.neosync_job.my_job.id
}
