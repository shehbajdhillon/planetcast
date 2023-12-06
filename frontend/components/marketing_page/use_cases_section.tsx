import {
  AspectRatio,
  Box,
  Button,
  HStack,
  Heading,
  Stack,
  Text,
  useColorModeValue
} from "@chakra-ui/react";
import { useState } from "react";

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


interface InfoGridViewProps extends InfoViewProps, VideoViewProps {
  flip: boolean;
}

const InfoGridView: React.FC<InfoGridViewProps> = (props) => {

  const { flip, headings, subheadings, transformations } = props;

  return (
    <Stack
      direction={{ base: "column", md: !flip ? "row" : "row-reverse" }}
      alignItems={"center"}
      spacing="25px"
    >
      <Box w="full" maxW={{ md: "60%" }}>
        <InfoView
          headings={headings}
          subheadings={subheadings}
        />
      </Box>
      <Box w="full" maxW={{ md: "40%" }}>
        <VideoView transformations={transformations} />
      </Box>
    </Stack>
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
      <AspectRatio ratio={16/9}>
        <Box display={"flex"} h="full" w="full" rounded={"sm"}>
          <video src={transformations[tfnIdx].link} controls />
        </Box>
      </AspectRatio>
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


const UseCasesSection: React.FC = () => {

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

  const headings2 = ["Marketing & Journalism"];
  const subheadings2 = [
    "Increase the reach of your breaking stories",
    "Ensure everyone stays informed with the latest trends"
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
      >
        Dubbing for your {' '}
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
      >
        Podcasts, training & educational videos, marketing content, and journalism media
      </Heading>

      <Stack w="full" spacing={{ base: "100px", md: "150px" }} pt={{ base:"50px", md: "100px" }}>

        <InfoGridView
          headings={headings1}
          subheadings={subheadings1}
          transformations={transformations1}
          flip={false}
        />

        <InfoGridView
          headings={headings2}
          subheadings={subheadings2}
          transformations={transformations2}
          flip={true}
        />

        <InfoGridView
          headings={headings3}
          subheadings={subheadings3}
          transformations={transformations3}
          flip={false}
        />
      </Stack>

    </Stack>
  );
};

export default UseCasesSection;
