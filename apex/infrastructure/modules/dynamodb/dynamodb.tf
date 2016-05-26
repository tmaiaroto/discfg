resource "aws_dynamodb_table" "DiscfgTable" {
    name = "${var.dynamodb_discfg_table}"
    read_capacity = 2
    write_capacity = 2
    hash_key = "key"
    attribute {
      name = "key"
      type = "S"
    }
}