import {
  Badge,
  Box,
  Button,
  HStack,
  Heading,
  IconButton,
  useColorMode,
  useColorModeValue,
} from '@chakra-ui/react';
import { useUser } from '@clerk/nextjs';
import { Moon, Sun } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';

interface NavbarProps {
  marketing?: boolean;
};

const Navbar: React.FC<NavbarProps> = ({ marketing }) => {

  const { toggleColorMode } = useColorMode();
  const { isSignedIn, isLoaded } = useUser();

  return (
    <Box w="full" display={"flex"} alignItems={"center"} justifyContent={"center"}>
      <Box
        display={"flex"}
        justifyContent={"space-between"}
        w="full"
        background={useColorModeValue("white", "black")}
        maxW={"1920px"}
      >
        <Box display={"flex"} alignItems={"center"} justifyContent={"center"}>
          <Image
            src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
            width={60}
            height={100}
            alt='spend sense logo'
          />
          <Heading
            fontSize={"30px"}
            display={{ base: "none", md:"flex" }}
            fontWeight={"semibold"}
          >
            PlanetCast
          </Heading>
          <Badge
            ml="5px"
          >
            Beta
          </Badge>
        </Box>
        <HStack>
          <IconButton
            onClick={toggleColorMode}
            aria-label='color mode toggle'
            icon={useColorModeValue(<Moon />, <Sun />)}
            variant={"ghost"}
          />
          <Link
            href={'/dashboard'}
            hidden={!(marketing && isLoaded)}
          >
            <Button
              backgroundColor={useColorModeValue("black", "white")}
              textColor={useColorModeValue("white", "black")}
              borderWidth={"1px"}
              _hover={{
                backgroundColor: useColorModeValue("white", "black"),
                textColor: useColorModeValue("black", "white")
              }}
            >
              { isSignedIn ? 'Dashboard' : 'Log In' }
            </Button>
          </Link>
        </HStack>
      </Box>
    </Box>
  );
}

export default Navbar;
