import React, { useRef, useState } from "react";
import axios, { AxiosError, AxiosProgressEvent } from "axios";
import { useAuth } from "../hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";

const FileUpload: React.FC<{ OnUploadComplete: () => void }> = ({
  OnUploadComplete,
}) => {
  const { getAuthToken } = useAuth();
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      setError(null);
    }
  };

  const handleBrowseClick = () => {
    fileInputRef.current?.click();
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    try {
      const token = await getAuthToken();

      // Get pre signed url
      console.log("file type: ", selectedFile.type);
      const response = await axios.post(
        "https://4j1h7lzpf5.execute-api.us-east-2.amazonaws.com/dev/upload-url",
        {
          fileName: selectedFile.name,
          fileType: selectedFile.type || "application/octet-stream",
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      const { uploadUrl, fileID } = response.data;

      // upload file to s3
      await axios.put(uploadUrl, selectedFile, {
        headers: { "Content-Type": selectedFile.type },
        onUploadProgress: (progressEvent: AxiosProgressEvent) => {
          const percentCompleted = progressEvent.total
            ? Math.round((progressEvent.loaded * 100) / progressEvent.total)
            : 0;
          setUploadProgress(percentCompleted);
        },
      });

      setUploadProgress(0);
      setSelectedFile(null);
      if (fileInputRef.current) fileInputRef.current.value = "";
      OnUploadComplete();
    } catch (error) {
      console.error("Upload failed:", error);
      setUploadProgress(0);
      if (axios.isAxiosError(error)) {
        const axiosError = error as AxiosError;
        setError(`Upload failed: ${axiosError.message}`);
      } else {
        setError("Upload failed: An unexpected error occurred");
      }
    }
  };

  return (
    <Card className="mt-4 border-dashed border-2 border-gray-300">
      <CardContent className="flex flex-col items-center justify-center py-10">
        <div className="rounded-full bg-gray-100 p-3 mb-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-gray-500"
          >
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="12" y1="18" x2="12" y2="12"></line>
            <line x1="9" y1="15" x2="15" y2="15"></line>
          </svg>
        </div>
        <p className="text-sm text-gray-500 mb-2">
          Select or Drag & Drop your files for upload
        </p>
        <p className="text-xs text-gray-400 mb-4">File size limit: 2 GB</p>
        <Input
          type="file"
          ref={fileInputRef}
          onChange={handleFileSelect}
          className="hidden"
          id="file-upload"
        />
        <label htmlFor="file-upload">
          <span className="inline-block">
            <Button variant="outline" onClick={handleBrowseClick}>
              Browse files...
            </Button>
          </span>
        </label>
        {selectedFile && (
          <div className="mt-4 text-center">
            <p>Selected file: {selectedFile.name}</p>
            <p>Size: {selectedFile.size} bytes</p>
            <Button onClick={handleUpload} className="mt-2">
              Upload
            </Button>
          </div>
        )}
        {uploadProgress > 0 && (
          <div className="w-full mt-4">
            <Progress value={uploadProgress} className="w-full" />
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default FileUpload;
