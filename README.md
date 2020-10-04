CIDR4
============

A simple library to make calculations on a set of IP address.


## Example

1. Install cidr4calc with `go get`

   ```
   go get haruue.moe/x/cidr4/cmd/cidr4calc
   ```

2. Create a script file `cidr4set.txt`

   ```
   # Empty lines and lines start with # will be omitted

   # CIDRs start with a PLUS(+) sign will make the CIDR plus into the CIDR set
   +0.0.0.0/1
   # The PLUS(+) sign can also be omitted
   128.0.0.0/1
   # CIDRs start with a MINUS(-) sign will make the CIDR minus from the CIDR set
   -192.0.2.0/24
   ```

3. Run `cidr4calc` with pipe

   ```
   $ cat cidr4set.txt | cidr4calc
   0.0.0.0/1
   128.0.0.0/2
   192.0.0.0/23
   192.0.3.0/24
   192.0.4.0/22
   192.0.8.0/21
   ...
   ```

## Library Usage

see [`ExampleCIDRv4Set_Minus`](cidrset_test.go#L13) for an example to use this
library.


