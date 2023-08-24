package hash

import (
	"fmt"
	"strings"
)

/*
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Formatter;

public class HashGenerator {

	public static String generateFixedHash(String input) {
	    String hash = sha256String(input.getBytes(StandardCharsets.UTF_8));
	    return String.format("%s-%s-%s-%s", hash.substring(0, 5), hash.substring(5, 10), hash.substring(10, 14), hash.substring(14, 17)).toUpperCase();
	}

	private static String sha256String(byte[] input) {
	    try {
	        MessageDigest digest = MessageDigest.getInstance("SHA-256");
	        byte[] hash = digest.digest(input);
	        return byteArrayToHexString(hash);
	    } catch (NoSuchAlgorithmException e) {
	        throw new RuntimeException("SHA-256 algorithm not available", e);
	    }
	}

	private static String byteArrayToHexString(byte[] bytes) {
	    try (Formatter formatter = new Formatter()) {
	        for (byte b : bytes) {
	            formatter.format("%02x", b);
	        }
	        return formatter.toString();
	    }
	}

	public static void main(String[] args) {
			String input = "your_input_string_here";
			String fixedHash = generateFixedHash(input);
			System.out.println(fixedHash);
		}
	}
*/
func GenerateFixedHash(input string) string {
	hash := sha256String([]byte(input))
	return strings.ToUpper(fmt.Sprintf("%s-%s-%s-%s", hash[0:5], hash[5:10], hash[10:14], hash[14:17]))
}
