#!/usr/bin/env ruby
# encoding: UTF-8

# exit 0 unless `git symbolic-ref HEAD`.end_with? "/master\n"

workspace = ENV["GOPATH"].split(":")[0]
Dir.chdir(workspace << "/src/github.com/registrobr/rdap")

final_status = 0

output = `go test -cover -race ./...`
lines = output.split("\n")

if $?.exitstatus != 0 then
    puts "Testes de unidade"
    lines.each {|line| puts line unless line.slice(0..1) =~ /^(ok|\?)/}
    final_status = 1
    puts "\n"
end

lines.delete_if do |line|
    coverage = line.scan(/(\d{1,3}.\d)%/)
    coverage.empty? || coverage[0][0].to_f >= 80
end

if !lines.empty? then
    puts "Os seguintes pacotes têm menos de 80% de cobertura:"
    lines.each {|line| puts "\t" << line.split("\t").drop(1).join("\t")}
    final_status = 1
end

def check_status(name, command)
    output = `#{command}`
    puts name << "\n\t" << output unless output.empty?
    return $?.exitstatus
end

def check_output(name, command)
    output = `#{command}`
    puts name << "\n\t" << output unless output.empty?
    return output.empty? ? 0 : 1
end

final_status += check_status("go vet", "go vet ./...")
final_status += check_status("defercheck", "defercheck ./...")
final_status += check_output("golint", "golint ./... | grep -v comment")
final_status += check_output("go-nyet", "go-nyet ./... | grep -v Godeps")
final_status += check_output("varcheck", "varcheck ./... | egrep -v ':( type| status)'")
final_status += check_output("gocyclo", "gocyclo -over 15 . | egrep -v '( Godeps|_test\.go| tests)'")

exit final_status

