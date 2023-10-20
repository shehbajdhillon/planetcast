import {
  Box,
  Button,
  HStack,
  Heading,
  IconButton,
  VStack,
  useColorMode,
  useColorModeValue,
  useDisclosure,
} from '@chakra-ui/react';
import { useUser } from '@clerk/nextjs';
import { Moon, Sun, X, Menu, XIcon } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useEffect } from 'react';
import NProgress from 'nprogress';

const Links = [
  {
    title: 'Home',
    link: '/#',
  },
  {
    title: 'Benefits',
    link: '/#usecases',
  },
  {
    title: 'Testimonials',
    link: '/#testimonials',
  },
  {
    title: 'Pricing',
    link: '/#pricing',
   },
  // {
  //   title: 'Blog',
  //   link: '/blog',
  // },
  //
];

export const NavLink = ({ children }: { children: React.ReactNode }) => (
  <Button
    px={2}
    py={1}
    rounded={'md'}
    variant={"ghost"}
    _hover={{
      backgroundColor: useColorModeValue("black", "white"),
      textColor: useColorModeValue("white", "black"),
    }}
  >
    {children}
  </Button>
);

interface NavbarProps {
  marketing?: boolean;
};

const Navbar: React.FC<NavbarProps> = ({ marketing }) => {

  const { toggleColorMode } = useColorMode();
  const { isSignedIn, isLoaded } = useUser();

  const router = useRouter();

  useEffect(() => {
    const handleRouteStart = () => NProgress.start();
    const handleRouteDone = () => NProgress.done();

    router.events.on('routeChangeStart', handleRouteStart);
    router.events.on('routeChangeComplete', handleRouteDone);
    router.events.on('routeChangeError', handleRouteDone);

    return () => {
      handleRouteDone();
      router.events.off('routeChangeStart', handleRouteStart);
      router.events.off('routeChangeComplete', handleRouteDone);
      router.events.off('routeChangeError', handleRouteDone);
    };
  }, [router.events]);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const navLinkColor = useColorModeValue('black', 'white');


  return (
    <Box
      w="full"
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      flexDir={"column"}
    >

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
            alt='planet cast logo'
          />
          <Heading
            fontSize={"30px"}
            display={{ base: "none", md:"flex" }}
            fontWeight={"medium"}
          >
            PlanetCast
          </Heading>
        </Box>

        <HStack>
          <HStack
            as={'nav'}
            spacing={4}
            display={{ base: 'none', lg: 'flex' }}
            marginRight={'15px'}
          >
            {Links.map((link) => (
              <Link href={link.link} key={link.title}>
                <NavLink>{link.title}</NavLink>
              </Link>
            ))}
          </HStack>

          <IconButton
            onClick={toggleColorMode}
            aria-label='color mode toggle'
            icon={useColorModeValue(<Moon />, <Sun />)}
            variant={"ghost"}
          />

          <IconButton
            onClick={isOpen ? onClose : onOpen}
            aria-label={'Open Menu'}
            icon={ isOpen ? <X /> : <Menu /> }
            display={{ base: 'inherit', lg: 'none' }}
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

      {isOpen ? (
        <Box pb={4} display={{ lg: 'none' }} w="full">
          <VStack
            as={'nav'}
            spacing={4}
            textColor={navLinkColor}
            fontWeight="600"
            alignItems={'left'}
            marginLeft={'25px'}
            marginTop={'20px'}
          >
            {Links.map((link) => (
              <Link href={link.link} key={link.title} onClick={onClose}>
                <NavLink>{link.title}</NavLink>
              </Link>
            ))}
          </VStack>
        </Box>
      ) : null}

    </Box>
  );
}

export default Navbar;
