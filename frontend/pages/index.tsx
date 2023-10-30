import Navbar from '@/components/marketing_page/marketing_page_navbar'
import {
  Box,
  HStack,
  Stack,
  Heading,
  Button,
  Text,
  useColorModeValue,
  VStack,
  useBreakpointValue,
  Grid,
  GridItem,
  Avatar,
  Center,
} from '@chakra-ui/react';
import Head from 'next/head'

import Image from 'next/image';
import Link from 'next/link';
import {
  ArrowUpFromDot,
  CircleDollarSign,
  DollarSign,
  ExternalLink,
  Globe2Icon,
  GlobeIcon,
  LanguagesIcon,
  TimerReset,
  TrendingDownIcon,
  Volume2Icon,
  ImportIcon,
} from 'lucide-react';
import VideoPlayer from '@/components/video_player';
import { useState } from 'react';
import PricingComponent from '@/components/marketing_page/pricing_component';
import FooterComponent from '@/components/marketing_page/footer_component';


const HeroSection: React.FC = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      direction={{ base: "column-reverse", md: "row" }}
      maxW={"1400px"}
      w="full"
    >
      <Box
        mb={{ base: "auto", md: "0px" }}
        w="full"
        maxW={{ md: "75%" }}
        alignItems={{ base: "center", md: "left" }}
        justifyContent={{ base: "center", md: "left" }}
        display={"flex"}
        flexDir={"column"}
      >
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Dub
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Translate
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Broadcast
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Content Across the {' '}
          <Text
            as={"span"}
            bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
            bgClip='text'
          >
            Planet
          </Text>
        </Heading>
        <HStack w={{ md: "full" }} pt="10px">
          <Link href={'/dashboard'}>
            <Button
              size={"lg"}
              backgroundColor={useColorModeValue("black", "white")}
              textColor={useColorModeValue("white", "black")}
              borderColor={useColorModeValue("white", "black")}
              borderWidth={"1px"}
              _hover={{
                backgroundColor: useColorModeValue("white", "black"),
                textColor: useColorModeValue("white", "black"),
                bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
              }}
            >
              Start for Free
            </Button>
          </Link>
          <Link href={"#usecases"}>
            <Button
              size={"lg"}
              variant={"outline"}
              _hover={{
                backgroundColor: useColorModeValue("black", "white:"),
                textColor: useColorModeValue("white", "black"),
                bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
              }}
            >
              Read More
            </Button>
          </Link>
        </HStack>
      </Box>
      <Box maxW={{ base: "200px", md: "25%" }} mt={{ base:"auto", md: "0px" }}>
        <Image
          height={1000}
          width={400}
          src={useColorModeValue('/planetcastgradientlight.svg', '/planetcastgradientdark.svg')}
          alt='planet cast gradient logo'
        />
      </Box>
    </Stack>
  );
};


const BenefitsSection: React.FC = () => {
  const iconSize = useBreakpointValue({ base: '40px', md: '60px' })
  const iconSizeBig = useBreakpointValue({ base: "75px", sm: "120px", md: "140px" });

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >

      <Heading
        size={{ base: '2xl', md: '3xl' }}
        fontWeight={'medium'}
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        <HStack>
          <Text>10x your {' '}
            <Text
              as={"span"}
              bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
              bgClip='text'
            >
              reach
            </Text>
          </Text>
          <ArrowUpFromDot size={iconSize} />
          <GlobeIcon size={iconSize} />
        </HStack>
        <HStack mt={{ md: "8px" }}>
          <Text>1/10th the {' '}
            <Text
              as={"span"}
              bgGradient={'linear(to-tr, #01CF00, #90DD00)'}
              bgClip='text'
            >
              cost
            </Text>
          </Text>
          <TrendingDownIcon size={iconSize} />
          <DollarSign size={iconSize} />
        </HStack>
      </Heading>

      <Box
        mb={{ base: "auto", md: "0px" }}
        w="full"
        alignItems={"left"}
        justifyContent={"left"}
        display={"flex"}
        flexDir={"column"}
      >
        <Stack direction={{ base: "column", md: "row" }} mt={{ base: "25px" }} spacing={{ base: "25px", md: "0px" }}>

          <Stack w="full" direction={{ base: "row", md: "column" }}>
            <Heading
              fontWeight={'normal'}
              textAlign={"left"}
              w={"full"}
              size={{ base: "sm", sm: "lg" }}
              px={{ base: "15px", md: "0px" }}
            >
              Engage listeners everywhere around the globe
            </Heading>
            <HStack w="full" px={{ base: "15px", md: "0px"}}>
              <Globe2Icon size={iconSizeBig} strokeWidth={"0.75px"}/>
              <LanguagesIcon size={iconSizeBig} strokeWidth={"0.75px"} />
            </HStack>
          </Stack>

          <Stack w="full" direction={{ base: "row-reverse", md: "column" }}>
            <Heading
              textAlign={"left"}
              w={"full"}
              fontWeight={'normal'}
              size={{ base: "sm", sm: "lg" }}
              px={{ base: "15px", md: "0px" }}
            >
              Save time and money over traditional dubbing
            </Heading>
            <HStack w="full" px={{ base: "15px", md: "0px"}}>
              <CircleDollarSign size={iconSizeBig} strokeWidth={"0.75px"} />
              <TimerReset size={iconSizeBig} strokeWidth={"0.75px"} />
            </HStack>
          </Stack>

          <Stack w="full" direction={{ base: "row", md: "column" }}>
            <Heading
              textAlign={"left"}
              w={"full"}
              fontWeight={'normal'}
              size={{ base: "sm", sm: "lg" }}
              px={{ base: "15px", md: "0px" }}
            >
              Preserve original voices in every translation
            </Heading>
            <HStack w="full" px={{ base: "15px", md: "0px"}}>
              <Volume2Icon size={iconSizeBig} strokeWidth={"0.75px"} />
              <ImportIcon size={iconSizeBig} strokeWidth={"0.75px"} />
            </HStack>
          </Stack>

        </Stack>
      </Box>
    </Stack>
  );
};


interface TestimonialCardProps {
  name: string;
  title: string;
  src: string;
  text: string[];
  link: string;
};

const TestimonialCard: React.FC<TestimonialCardProps> = (props) => {
  const { name, title, src, text, link } = props;

  return (
    <Box
      flex="1"
      p={6}
      shadow="lg"
      borderRadius="md"
      borderWidth={1}
      position="relative"
    >
      <HStack mb="30px">
        <Avatar src={src} name={name} />
        <Stack spacing={1}>
          <HStack>
            <Heading
              size={{ base: "md" }}
              fontWeight={'medium'}
            >
              { name }
            </Heading>
            <Link
              href={link}
              target='_blank'
              aria-label={`Read more about ${name}`}
            >
              <ExternalLink />
            </Link>
          </HStack>
          <Heading
            size={{ base:"sm" }}
            fontWeight={'normal'}
          >
            { title }
          </Heading>
        </Stack>
      </HStack>

      <Heading
        fontSize={"md"}
        fontWeight={"normal"}
        as={"em"}
      >
        {text.map((txt, idx) => (
          <Text pt="10px" key={idx}>{'"'}{txt}{'"'}</Text>
        ))}
      </Heading>
    </Box>
  );
};


const TestimonialSection = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        Welcome to efficient {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          broadcasting
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        Discover how our users are revolutionizing their content reach
      </Heading>

      <Stack
        direction={{ base: "column", md: "row" }}
        justifyContent={{ md: "space-between" }}
        w="full"
        mt="45px"
        px="16px"
        spacing={"80px"}
      >
        <TestimonialCard
          name='Rahul Pandey'
          title='Educator'
          src={'/rahulimg.jpeg'}
          text={['Amazed by this 🤯😮', 'The next few years for creators are going to be wild.']}
          link='https://www.youtube.com/@RahulInHindi'
        />
        <TestimonialCard
          name='Devin Estopinal'
          title='Social Media Content Strategist'
          src={'/devnimg.jpeg'}
          text={["I can't wait to see the results", "Just got done downloading the video after dubbing in spanish and it is absolutely flawless"]}
          link='https://x.com/NotDevn'
        />
        <TestimonialCard
          name='Dallon Asnes'
          title='Travel Content Creator'
          src={'/dallonimg.jpeg'}
          text={['first vid in progress', 'the hindi dubbing is excellent']}
          link='https://www.youtube.com/@dallonearth'
        />
      </Stack>
    </Stack>
  );
};


const UseCasesSection = () => {

  const headings1 = ["Training & Education"];
  const subheadings1 = [
    "Make your educational content more effective",
    "Employees and students can now understand training materials in their own tongue",
  ];
  const transformations1 = [
    {
      language: "ENGLISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/training_video_english.mp4",
    },
    {
      language: "SPANISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/training_video_spanish.mp4",
    },
    {
      language: "HINDI",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/training_video_hindi.mp4",
    },
    {
      language: "FRENCH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/training_video_french.mp4",
    },
  ]

  const headings2 = ["Journalism"];
  const subheadings2 = [
    "Increase the reach of your breaking news stories",
    "Ensure every household stays informed with the most current events"
  ];
  const transformations2 = [
    {
      language: "ENGLISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/journalism_english.mp4",
    },
    {
      language: "SPANISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/journalism_spanish.mp4",
    },
    {
      language: "HINDI",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/journalism_hindi.mp4",
    },
    {
      language: "FRENCH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/journalism_french.mp4",
    },
  ]

  const headings3 = ["Postcasts"];
  const subheadings3 = [
    "Amplify your podcast's resonance",
    "Connect with listeners worldwide by sharing episodes in their preferred language"
  ];
  const transformations3 = [
    {
      language: "ENGLISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/podcast_english.mp4",
    },
    {
      language: "SPANISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/podcast_spanish.mp4",
    },
    {
      language: "HINDI",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/podcast_hindi.mp4",
    },
    {
      language: "FRENCH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/podcast_french.mp4",
    },
  ]

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        Tailored for your {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          use cases
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        Whether you produce training & educational content, podcasts, or journalism media
      </Heading>

      <Stack w="full" spacing={{ base: "100px", md: "150px" }} pt={{ base:"50px", md: "100px" }}>

        <InfoGridView
          headings={headings1}
          subheadings={subheadings1}
          transformations={transformations1}
        />

        <InfoGridView
          headings={headings2}
          subheadings={subheadings2}
          transformations={transformations2}
          flip
        />

        <InfoGridView
          headings={headings3}
          subheadings={subheadings3}
          transformations={transformations3}
        />
      </Stack>

    </Stack>
  );
};


interface InfoGridViewProps extends InfoViewProps, VideoViewProps {
  flip?: boolean;
}

const InfoGridView: React.FC<InfoGridViewProps> = (props) => {

  const { flip, headings, subheadings, transformations } = props;

  return (
    <Stack
      direction={{ base: "column", md: !flip ? "row" : "row-reverse" }}
      alignItems={"center"}
      spacing="25px"
    >
      <Box w="full" px={{ base: "15px", md: "0px" }} maxW={{ md: "60%" }}>
        <InfoView
          headings={headings}
          subheadings={subheadings}
        />
      </Box>
      <Box w="full" px={{ base: "15px", md: "0px" }} maxW={{ md: "40%" }}>
        <VideoView transformations={transformations} />
      </Box>
    </Stack>
  );
};


interface InfoViewProps {
  headings: string[];
  subheadings: string[];
};

const InfoView: React.FC<InfoViewProps> = ({ headings, subheadings }) => {
  return (
    <Box
      mb={{ base: "auto", md: "0px" }}
      w="full"
      alignItems={{ base: "center", md: "left" }}
      justifyContent={{ base: "center", md: "left" }}
      display={"flex"}
      flexDir={"column"}
    >
      <Heading
        size={{ base: '2xl', md: '3xl' }}
        fontWeight={'medium'}
        mb={{ md: "10px" }}
        textAlign={"left"}
        w={"full"}
      >
        {headings.map((heading, idx) => (
          <HStack key={idx}>
            <Text key={idx}>{heading}</Text>
          </HStack>
        ))}
      </Heading>
      {subheadings.map((heading, idx) => (
        <Heading
          fontWeight={'normal'}
          size={{ base: "sm", sm: "lg" }}
          key={idx}
          textAlign={"left"}
          w={"full"}
        >
          {heading}
        </Heading>
      ))}
    </Box>
  );
};


interface VideoViewProps {
  transformations: Record<string, any>[];
};

const VideoView: React.FC<VideoViewProps> = ({ transformations }) => {

  const buttonBg = useColorModeValue("black", "white");
  const buttonColor = useColorModeValue("white", "black");
  const [tfnIdx, setTfnIdx] = useState(0);

  return (
    <Box w="full">
      <Box display={"flex"} h="full" w="full" rounded={"sm"}>
        <VideoPlayer src={transformations[tfnIdx].link} />
      </Box>
      <HStack pt="10px">
        {transformations.map((tfn, idx) => (
          <Button
            key={idx}
            onClick={() => setTfnIdx(idx)}
            variant={idx == tfnIdx ? "solid" : "outline"}
            pointerEvents={idx === tfnIdx ? "none" : "auto"}
            background={idx === tfnIdx ? buttonBg : buttonColor }
            color={idx === tfnIdx ? buttonColor : '' }
          >
            {tfn.language}
          </Button>
        ))}
      </HStack>
    </Box>
  );
};


const PricingSection = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        Start dubbing {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          today
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        mb="45px"
        textAlign={"left"}
        w={"full"}
        px={{ base: "15px", md: "0px" }}
      >
        Select the perfect plan tailored to your needs
      </Heading>
      <PricingComponent />
    </Stack>
  );
};


export default function Home() {
  const bgColor = useColorModeValue("white", "black");
  return (
    <VStack>
      <Head>
        <title>PlanetCast</title>
        <meta name="description" content="Cast your Content Across the Planet" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box position={"fixed"} top={0} left={0} w="full" px="10px" backgroundColor={bgColor} zIndex={100}>
        <Navbar marketing />
      </Box>
      <VStack w="full">

        <Center
          w="full"
          py={{ base:"110px", md: "250px" }}
        >
          <HeroSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="usecases"
        >
          <UseCasesSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="benefits"
        >
          <BenefitsSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="testimonials"
        >
          <TestimonialSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="pricing"
        >
          <PricingSection />
        </Center>

        <FooterComponent />

      </VStack>
    </VStack>
  )
}

