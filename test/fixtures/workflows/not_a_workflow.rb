puts "Hello, world!"
input = STDIN.gets.chomp
def permutations(str)
  return [str] if str.length <= 1
  perms = []
  str.chars.each_with_index do |char, idx|
    rest = str[0...idx] + str[(idx+1)..-1]
    permutations(rest).each do |perm|
      perms << char + perm
    end
  end
  perms.uniq
end

permutations(input).each { |perm| puts perm }
