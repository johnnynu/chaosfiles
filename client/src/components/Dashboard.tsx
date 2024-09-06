import FileUpload from "./FileUpload"; // Import the FileUpload component
import Navbar from "./Navbar";

const Dashboard = () => {
  return (
    <div className="container mx-auto p-4">
      <Navbar />
      <h1 className="text-3xl font-bold mb-8">Welcome!</h1>
      <p className="text-gray-600 mb-4">
        Share your files with ChaosFiles.com. Upload as much data as you want,
        download without speed limits and share files with friends without
        additional costs!
      </p>
      <FileUpload OnUploadComplete={() => {}} />
    </div>
  );
};

export default Dashboard;
